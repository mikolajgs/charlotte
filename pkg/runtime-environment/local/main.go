package localruntimeenvironment

import (
	runtimeenvironment "charlotte/pkg/runtime-environment"
	"charlotte/pkg/step"
	"fmt"
	"os"
	"os/exec"
)

type LocalRuntimeEnvironment struct {
	runtimeenvironment.RuntimeEnvironment
}

func (e *LocalRuntimeEnvironment) Create() error {
	return nil
}

func (c *LocalRuntimeEnvironment) Destroy() error {
	return nil
}

// Run runs a Step and returns error code, error string, path to file containing stdout and path to file containing stderr.
func (e *LocalRuntimeEnvironment) Run(step step.IStep) (string, string, error) {
	if _, err := exec.LookPath("bash"); err != nil {
		return "", "", fmt.Errorf("command bash not found: %w", err)
	}

	// fOut and fErr are io.File to stdout and stderr
	fStep, fOut, fErr, err := e.InitRunStep(step)
	if err != nil {
		return "", "", fmt.Errorf("error initializing step: %w", err)
	}

	defer os.Remove(fStep.Name())
	defer fStep.Close()
	defer fOut.Close()
	defer fErr.Close()

	cmd, stdout, stderr, err := e.CreateCmd("bash", fStep.Name())
	if err != nil {
		return "", "", fmt.Errorf("error creating command: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return "", "", fmt.Errorf("error starting command: %w", err)
	}

	// create wait group that attaches stdout and stderr to files
	wg := e.CreateWaitGroup(stdout, fOut, stderr, fErr)
	wg.Wait()

	// wait for the command to finish
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return fOut.Name(), fErr.Name(), fmt.Errorf("command returns exit code %s", exiterr)
		} else {
			return "", "", fmt.Errorf("error waiting for the command: %w", err)
		}
	}

	return fOut.Name(), fErr.Name(), nil
}
