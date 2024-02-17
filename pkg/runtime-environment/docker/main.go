package dockerruntimeenvironment

import (
	runtimeenvironment "charlotte/pkg/runtime-environment"
	"charlotte/pkg/step"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerRuntimeEnvironment struct {
	runtimeenvironment.RuntimeEnvironment
}

func (e *DockerRuntimeEnvironment) Create() error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error creating docker client: %w", err)
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("error pulling docker image: %w", err)
	}
	io.Copy(os.Stdout, reader)

	_, err = cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
	}, nil, nil, nil, "")
	if err != nil {
		return fmt.Errorf("error creating docker container: %w", err)
	}

	return nil
}

func (e *DockerRuntimeEnvironment) Destroy() error {
	return nil
}

/*func (e *DockerRuntimeEnvironment) CreateDocker(fStep string, fOut string, fErr string) error {
	return nil
}*/

/*cmd := exec.Command("docker", "container", "create", "-t", "-i", "alpine")
if err := cmd.Start(); err != nil {
	return "", fmt.Errorf(fmt.Sprint("error starting cmd: %s", err.Error()))
}
if err := cmd.Wait(); err != nil {
	if exiterr, ok := err.(*exec.ExitError); ok {
		return "", fmt.Errorf(fmt.Sprint("exit code different than 0: %d", exiterr.ExitCode()))
	} else {
		return "", fmt.Errorf(fmt.Sprint("error waiting for cmd: %s", err.Error()))
	}
}*/

/*return "id", nil
}*/

func (e *DockerRuntimeEnvironment) Run(step step.IStep) (string, string, error) {
	return "", "", nil
	/*errCode, errStr, fStep, fOut, fErr := e.InitRunStep(step)
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

	return 0, []string{}, fOut.Name(), fErr.Name()*/
}
