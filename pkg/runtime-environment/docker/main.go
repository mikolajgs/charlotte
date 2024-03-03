package dockerruntimeenvironment

import (
	runtimeenvironment "charlotte/pkg/runtime-environment"
	"charlotte/pkg/step"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type DockerRuntimeEnvironment struct {
	runtimeenvironment.RuntimeEnvironment
	containerId string
	client      *client.Client
}

func (e *DockerRuntimeEnvironment) Create() error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error creating docker client: %w", err)
	}
	e.client = cli

	reader, err := cli.ImagePull(ctx, "debian:bookworm-slim", types.ImagePullOptions{Platform: "linux/amd64"})
	if err != nil {
		return fmt.Errorf("error pulling docker image: %w", err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "debian:bookworm-slim",
		Cmd:   []string{"sleep", "600"},
	}, nil, nil,
		&ocispec.Platform{
			Architecture: "amd64",
			OS:           "linux",
		},
		"")
	if err != nil {
		return fmt.Errorf("error creating docker: %w", err)
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("error starting docker: %w", err)
	}

	e.containerId = resp.ID
	return nil
}

func (e *DockerRuntimeEnvironment) Destroy() error {
	ctx := context.Background()
	err := e.client.ContainerRemove(ctx, e.containerId, container.RemoveOptions{
		Force:         true,
		RemoveLinks:   true,
		RemoveVolumes: true,
	})
	if err != nil {
		return fmt.Errorf("error removing docker: %w", err)
	}

	if e.client != nil {
		e.client.Close()
	}
	return nil
}

func (e *DockerRuntimeEnvironment) Run(step step.IStep) (string, string, error) {
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

	fStepArchive, err := archive.Tar(fStep.Name(), archive.Gzip)
	if err != nil {
		return "", "", fmt.Errorf("error making tar: %w", err)
	}

	err = e.client.CopyToContainer(ctx, e.containerId, "/tmp/", fStepArchive, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
	if err != nil {
		return "", "", fmt.Errorf("error copying step script to docker container: %w", err)
	}

	execConfig := types.ExecConfig{
		Cmd:          []string{"bash", fmt.Sprintf("/tmp/%s", filepath.Base(fStep.Name()))},
		AttachStdout: true,
		AttachStderr: true,
	}
	execID, err := e.client.ContainerExecCreate(ctx, e.containerId, execConfig)
	if err != nil {
		return "", "", fmt.Errorf("error running bash with step: %w", err)
	}

	resp, err := e.client.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return "", "", fmt.Errorf("error attaching to exec: %w", err)
	}
	defer resp.Close()

	_, err = stdcopy.StdCopy(fOut, fErr, resp.Reader)
	if err != nil {
		return "", "", fmt.Errorf("error getting stdout and stderr: %w", err)
	}

	inspectedExec, err := e.client.ContainerExecInspect(ctx, execID.ID)
	if err != nil {
		return "", "", fmt.Errorf("error inspecting exec: %w", err)
	}

	if inspectedExec.ExitCode != 0 {
		return fOut.Name(), fErr.Name(), fmt.Errorf("command returns exit code %d", inspectedExec.ExitCode)
	}

	return fOut.Name(), fErr.Name(), nil
}
