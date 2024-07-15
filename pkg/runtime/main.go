package runtime

import (
	"charlotte/pkg/step"
	"errors"
	"fmt"
	"os"
)

type IRuntime interface {
	Run(step step.IStep, stepNumber int, env *map[string]string) (string, string, error)
	Create(steps []step.IStep) error
	Destroy(steps []step.IStep) error
}

type Runtime struct {
}

// InitStepScriptOutputs creates temporary files to write the stdout and write the stderr.
func (r *Runtime) InitStepScriptOutputs(step step.IStep) (*os.File, *os.File, error) {
	fOut, err := os.CreateTemp("", "stepout.*.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("error creating tmp file for stdout: %w", err)
	}

	fErr, err := os.CreateTemp("", "steperr.*.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("error creating tmp file for stderr: %w", err)
	}

	return fOut, fErr, nil
}

// InitStepScriptContents creates temporary file to store the script.
func (r *Runtime) InitStepScriptContents(step step.IStep) (*os.File, error) {
	script := step.GetRunScript()
	if script == "" {
		return nil, errors.New("script is empty")
	}

	fStep, err := os.CreateTemp("", "step.*.sh")
	if err != nil {
		return nil, fmt.Errorf("error creating tmp file for script: %w", err)
	}

	if _, err := fStep.Write([]byte(script)); err != nil {
		fStep.Close()
		return nil, fmt.Errorf("error writing tmp file: %w", err)
	}

	return fStep, nil
}

// TODO: Tidy up Init* functions
