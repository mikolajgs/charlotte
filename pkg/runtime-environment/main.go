package runtimeenvironment

import (
	"charlotte/pkg/step"
	"errors"
	"fmt"
	"os"
)

type IRuntimeEnvironment interface {
	Run(step step.IStep) (string, string, error)
	Create() error
	Destroy() error
}

type RuntimeEnvironment struct {
}

// InitRunStep will create three temporary files to store the script, write the stdout and write the stderr.
func (r *RuntimeEnvironment) InitRunStep(step step.IStep) (*os.File, *os.File, *os.File, error) {
	script := step.GetScript()
	if script == "" {
		return nil, nil, nil, errors.New("script is empty")
	}

	fStep, err := os.CreateTemp("", "step.*.sh")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating tmp file for script: %w", err)
	}

	if _, err := fStep.Write([]byte(script)); err != nil {
		fStep.Close()
		return nil, nil, nil, fmt.Errorf("error writing tmp file: %w", err)
	}

	fOut, err := os.CreateTemp("", "stepout.*.txt")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating tmp file for stdout: %w", err)
	}

	fErr, err := os.CreateTemp("", "steperr.*.txt")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating tmp file for stderr: %w", err)
	}

	return fStep, fOut, fErr, nil
}
