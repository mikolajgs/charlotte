package kubernetesruntimeenvironment

import (
	"bytes"
	runtimeenvironment "charlotte/pkg/runtime-environment"
	"charlotte/pkg/step"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

const DEFAULT_NAMESPACE = "charlotte"
const DOCKER_IMAGE = "mikogs/charlotte:0.1.0"

type KubernetesRuntimeEnvironment struct {
	runtimeenvironment.RuntimeEnvironment
	namespace      string
	config         *rest.Config
	clientSet      *kubernetes.Clientset
	InCluster      bool
	KubeconfigPath string
}

func (e *KubernetesRuntimeEnvironment) Create(steps []step.IStep) error {
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

	// Create a ConfigMap with all the step scripts
	configMapData := map[string]string{}
	for i, step := range steps {
		configMapData[fmt.Sprintf("script-%d.sh", i)] = step.GetScript()
	}
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "charlotte-cm",
		},
		Data: configMapData,
	}

	_, err = e.clientSet.CoreV1().ConfigMaps(DEFAULT_NAMESPACE).Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create configmap: %w", err)
	}

	// Create a pod with configmap attached to it
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "charlotte-pod",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    "charlotte-container",
					Image:   DOCKER_IMAGE,
					Command: []string{"sleep", "600"},
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      "charlotte-scripts",
							MountPath: "/tmp/charlotte-cm",
						},
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: "charlotte-scripts",
					VolumeSource: v1.VolumeSource{
						ConfigMap: &v1.ConfigMapVolumeSource{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "charlotte-cm",
							},
						},
					},
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

func (e *KubernetesRuntimeEnvironment) Destroy(steps []step.IStep) error {
	return nil
}

func (e *KubernetesRuntimeEnvironment) Run(step step.IStep, stepNumber int) (string, string, error) {
	// fOut and fErr are io.File to stdout and stderr
	fOut, fErr, err := e.InitStepOutputs(step)
	if err != nil {
		return "", "", fmt.Errorf("error initializing step outputs: %w", err)
	}
	defer fOut.Close()
	defer fErr.Close()

	ctx := context.Background()

	command := []string{"bash", fmt.Sprintf("/tmp/charlotte-cm/script-%d.sh", stepNumber)}
	req := e.clientSet.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name("charlotte-pod").
		Namespace(DEFAULT_NAMESPACE).
		SubResource("exec").
		Param("command", "echo").
		Param("stdin", "false").
		Param("stdout", "true").
		Param("stderr", "true").
		VersionedParams(&v1.PodExecOptions{
			Command: command,
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
		}, metav1.ParameterCodec)

	fmt.Fprintf(os.Stdout, "%v", req.URL())

	executor, err := remotecommand.NewSPDYExecutor(e.config, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("error creating executor: %w", err)
	}

	var stdout, stderr bytes.Buffer
	err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  strings.NewReader(""),
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		fmt.Printf("STDOUT: %s\n", stdout.String())
		fmt.Printf("STDERR: %s\n", stderr.String())
		return "", "", fmt.Errorf("error executing command in pod: %w", err)
	}

	fmt.Printf("STDOUT: %s\n", stdout.String())
	fmt.Printf("STDERR: %s\n", stderr.String())

	return fOut.Name(), fErr.Name(), nil
}
