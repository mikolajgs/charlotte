package kubernetesruntimeenvironment

import (
	"archive/tar"
	"bytes"
	runtimeenvironment "charlotte/pkg/runtime-environment"
	"charlotte/pkg/step"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

const DEFAULT_NAMESPACE = "charlotte"

type KubernetesRuntimeEnvironment struct {
	runtimeenvironment.RuntimeEnvironment
	namespace      string
	config         *rest.Config
	clientSet      *kubernetes.Clientset
	InCluster      bool
	KubeconfigPath string
}

func (e *KubernetesRuntimeEnvironment) Create() error {
	var config *rest.Config
	var err error
	if e.InCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			return fmt.Errorf("error building incluster config: %w", err)
		}
	} else {
		if e.KubeconfigPath == "" {
			return fmt.Errorf("kubeconfig path is missing")
		}
		config, err = clientcmd.BuildConfigFromFlags("", e.KubeconfigPath)
		if err != nil {
			return fmt.Errorf("error reading kubeconfig: %w", err)
		}
	}
	e.config = config

	clientSet, err := kubernetes.NewForConfig(e.config)
	if err != nil {
		return fmt.Errorf("error building kubernetes clientset: %w", err)
	}
	e.clientSet = clientSet

	if e.namespace != "" {
		e.namespace = DEFAULT_NAMESPACE
	}

	// Get the list of namespaces
	namespaces, err := e.clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error getting namespaces: %w", err)
	}

	found := false
	for _, namespace := range namespaces.Items {
		if namespace.Name == DEFAULT_NAMESPACE {
			found = true
		}
	}

	if !found {
		namespace := &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: DEFAULT_NAMESPACE,
			},
		}
		_, err = e.clientSet.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("error creating namespace '%s': %w", DEFAULT_NAMESPACE, err)
		}
	}

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "charlotte-pod",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    "charlotte-container",
					Image:   "debian:bookworm-slim",
					Command: []string{"sleep", "600"},
				},
			},
		},
	}

	_, err = e.clientSet.CoreV1().Pods(DEFAULT_NAMESPACE).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create pod: %w", err)
	}

	time.Sleep(10 * time.Second)

	return nil
}

func (e *KubernetesRuntimeEnvironment) Destroy() error {
	return nil
}

func (e *KubernetesRuntimeEnvironment) Run(step step.IStep) (string, string, error) {
	// fOut and fErr are io.File to stdout and stderr
	fStep, fOut, fErr, err := e.InitRunStep(step)
	if err != nil {
		return "", "", fmt.Errorf("error initializing step: %w", err)
	}

	defer os.Remove(fStep.Name())
	defer fStep.Close()
	defer fOut.Close()
	defer fErr.Close()

	ctx := context.Background()

	filePath := fStep.Name()
	destPath := fmt.Sprintf("/tmp/%s", filepath.Base(fStep.Name()))
	fileContents, err := io.ReadAll(fStep)
	if err != nil {
		return "", "", fmt.Errorf("error reading step script: %w", err)
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", "", fmt.Errorf("error getting fileinfo for step script: %w", err)
	}

	tarFileName := fmt.Sprintf("%s.tar", fStep.Name())
	tarFile, err := os.Create(tarFileName)
	if err != nil {
		return "", "", fmt.Errorf("failed to create tar file: %w", err)
	}
	defer tarFile.Close()

	// Prepare the tar command to run in the pod
	command := []string{"tar", "xf", "-", "-C", destPath}
	req := e.clientSet.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name("charlotte-pod").
		Namespace(DEFAULT_NAMESPACE).
		SubResource("exec").
		Param("command", "echo").
		Param("command", "aaaa").
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		VersionedParams(&v1.PodExecOptions{
			Command: command,
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
		}, metav1.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(e.config, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("error creating executor: %w", err)
	}

	header := &tar.Header{
		Name: fileInfo.Name(),
		Mode: int64(fileInfo.Mode()),
		Size: fileInfo.Size(),
	}

	tw := tar.NewWriter(tarFile)
	defer tw.Close()

	if err := tw.WriteHeader(header); err != nil {
		return "", "", fmt.Errorf("error creating tar writer: %w", err)
	}

	if _, err := tw.Write(fileContents); err != nil {
		return "", "", fmt.Errorf("error creating tar archive: %w", err)
	}

	// Execute the command in the pod
	var stdout, stderr bytes.Buffer
	err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  tarFile,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return "", "", fmt.Errorf("error executing command in pod: %w", err)
	}

	fmt.Printf("STDOUT: %s\n", stdout.String())
	fmt.Printf("STDERR: %s\n", stderr.String())

	return fOut.Name(), fErr.Name(), nil
}
