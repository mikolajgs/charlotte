package kubernetesruntimeenvironment

import (
	runtimeenvironment "charlotte/pkg/runtime-environment"
	"charlotte/pkg/step"
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes/clientset"
	"k8s.io/client-go/rest"
)

type KubernetesRuntimeEnvironment struct {
	runtimeenvironment.RuntimeEnvironment
	namespace string
	config    *rest.Config
	clientSet *clientset.Clientset
}

func (e *KubernetesRuntimeEnvironment) Create() error {
	e.config, err = rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("error building incluster config: %w", err)
	}

	e.clientSet, err := clientset.NewForConfig(e.config)
	if err != nil {
		return fmt.Errorf("error building kubernetes clientset: %w", err)
	}

	// if namespace does not exist, it has to be created
	// and use some default name to it

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
