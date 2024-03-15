package kubernetesruntimeenvironment

import (
	runtimeenvironment "charlotte/pkg/runtime-environment"
	"charlotte/pkg/step"
	"context"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const DEFAULT_NAMESPACE = "charlotte"

type KubernetesRuntimeEnvironment struct {
	runtimeenvironment.RuntimeEnvironment
	namespace string
	config    *rest.Config
	clientSet *kubernetes.Clientset
}

func (e *KubernetesRuntimeEnvironment) Create() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("error building incluster config: %w", err)
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

	// Do something
	// if namespace does not exist

	return fOut.Name(), fErr.Name(), nil
}
