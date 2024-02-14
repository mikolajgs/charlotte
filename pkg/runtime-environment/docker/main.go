package dockerruntimeenvironment

import (
	runtimeenvironment "charlotte/pkg/runtime-environment"
	"charlotte/pkg/step"
	"fmt"
	"os"
	"os/exec"
)

type DockerRuntimeEnvironment struct {
	runtimeenvironment.RuntimeEnvironment
}

func (e *DockerRuntimeEnvironment) CreateDocker(fStep string, fOut string, fErr string) (string, error) {
	cmd := exec.Command("docker", "container", "create", "-t", "-i", "alpine")
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf(fmt.Sprint("error starting cmd: %s", err.Error()))
	}
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf(fmt.Sprint("exit code different than 0: %d", exiterr.ExitCode()))
		} else {
			return "", fmt.Errorf(fmt.Sprint("error waiting for cmd: %s", err.Error()))
		}
	}

	return "id", nil
}

func (e *DockerRuntimeEnvironment) Run(step step.IStep) (int, []string, string, string) {
	if _, err := exec.LookPath("docker"); err != nil {
		return 102, []string{fmt.Sprint("docker not found: %s", err.Error())}, "", ""
	}

	errCode, errStr, fStep, fOut, fErr := e.InitRunStep(step)
	if errCode != 0 {
		return errCode, []string{errStr}, "", ""
	}

	defer os.Remove(fStep.Name())
	defer fStep.Close()
	defer fOut.Close()
	defer fErr.Close()

	_, err := e.CreateDocker(fStep.Name(), fOut.Name(), fErr.Name())
	if err != nil {
		return 110, []string{fmt.Sprintf("error creating docker container: %s", err.Error())}, "", ""
	}

	errCode, errStr, cmd, stdout, stderr := e.CreateCmd("bash", fStep.Name())
	if errCode != 0 {
		return errCode, []string{errStr}, "", ""
	}

	if err := cmd.Start(); err != nil {
		return 108, []string{fmt.Sprint("error starting cmd: %s", err.Error())}, "", ""
	}

	wg := e.CreateWaitGroup(stdout, fOut, stderr, fErr)
	wg.Wait()

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return exiterr.ExitCode(), []string{"exit code different than 0"}, fOut.Name(), fErr.Name()
		} else {
			return 109, []string{fmt.Sprint("error waiting for cmd: %s", err.Error())}, "", ""
		}
	}

	return 0, []string{}, fOut.Name(), fErr.Name()
}
