package localruntimeenvironment

import (
	"fmt"
	"os"
	"charlotte/pkg/step"
	"os/exec"
	"charlotte/pkg/runtime-environment"
)

type LocalRuntimeEnvironment struct {
	runtimeenvironment.RuntimeEnvironment
}

// Run runs a Step and returns error code, error string, path to file containing stdout and path to file containing stderr.
func (e *LocalRuntimeEnvironment) Run(step step.IStep) (int, []string, string, string) {
	if _, err := exec.LookPath("bash"); err != nil {
		return 102, []string{fmt.Sprint("bash not found: %s", err.Error())}, "", ""
	}

	// fOut and fErr are io.File to stdout and stderr
	errCode, errStr, fStep, fOut, fErr := e.InitRunStep(step)
	if errCode != 0 {
		return errCode, []string{errStr}, "", ""
	}
	
	defer os.Remove(fStep.Name())
	defer fStep.Close()
	defer fOut.Close()
	defer fErr.Close()

	errCode, errStr, cmd, stdout, stderr := e.CreateCmd("bash", fStep.Name())
	if errCode != 0 {
		return errCode, []string{errStr}, "", ""
	}

	if err := cmd.Start(); err != nil {
		return 108, []string{fmt.Sprint("error starting cmd: %s", err.Error())}, "", ""
	}

	// create wait group that attaches stdout and stderr to files
	wg := e.CreateWaitGroup(stdout, fOut, stderr, fErr)
	wg.Wait()

	// wait for the command to finish
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return exiterr.ExitCode(), []string{"exit code different than 0"}, fOut.Name(), fErr.Name()
		} else {
			return 109, []string{fmt.Sprint("error waiting for cmd: %s", err.Error())}, "", ""
		}
	}

	return 0, []string{}, fOut.Name(), fErr.Name()
}
