package localruntime

import (
	"bufio"
	runtime "charlotte/pkg/runtime"
	"charlotte/pkg/step"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type LocalRuntime struct {
	runtime.Runtime
}

func (e *LocalRuntime) Create(steps []step.IStep) error {
	return nil
}

func (c *LocalRuntime) Destroy(steps []step.IStep) error {
	return nil
}

// Run runs a Step and returns error code, error string, path to file containing stdout and path to file containing stderr.
func (e *LocalRuntime) Run(step step.IStep, stepNumber int, env *map[string]string) (string, string, error) {
	if _, err := exec.LookPath("bash"); err != nil {
		return "", "", fmt.Errorf("command bash not found: %w", err)
	}

	// fStep is io.File with contents of the script
	fStep, err := e.InitStepScript(step)
	if err != nil {
		return "", "", fmt.Errorf("error initializing step script: %w", err)
	}

	// fOut and fErr are io.File to stdout and stderr
	fOut, fErr, err := e.InitStepOutputs(step)
	if err != nil {
		return "", "", fmt.Errorf("error initializing step outputs: %w", err)
	}

	defer os.Remove(fStep.Name())
	defer fStep.Close()
	defer fOut.Close()
	defer fErr.Close()

	cmd, stdout, stderr, err := e.CreateCmd(env, "bash", fStep.Name())
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

// CreateCmd creates and returns a command along with io.ReadCloser to attach stdout and stderr.
func (e *LocalRuntime) CreateCmd(env *map[string]string, name string, args ...string) (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	cmd := exec.Command(name, args...)

	if env != nil && len(*env)>0 {
		for k, v := range *env {
			cmd.Env = append(cmd.Environ(), fmt.Sprintf("%s=%s", k, v))
		}
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error piping stdout: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error piping stderr: %w", err)
	}

	return cmd, stdout, stderr, nil
}

// CreateWaitGroup creates WaitGroups that will pipe stdout and stderr to specific files.
func (e *LocalRuntime) CreateWaitGroup(stdout io.ReadCloser, fOut *os.File, stderr io.ReadCloser, fErr *os.File) *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		scanner := bufio.NewScanner(stdout)
		writer := io.MultiWriter(fOut, os.Stdout)
		for scanner.Scan() {
			fmt.Fprintln(writer, scanner.Text())
		}
		wg.Done()
	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		writer := io.MultiWriter(fErr, os.Stderr)
		for scanner.Scan() {
			fmt.Fprintln(writer, scanner.Text())
		}
		wg.Done()
	}()
	return wg
}
