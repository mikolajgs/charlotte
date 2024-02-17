package runtimeenvironment

import (
	"bufio"
	"charlotte/pkg/step"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
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

// CreateCmd creates and returns a command along with io.ReadCloser to attach stdout and stderr.
func (r *RuntimeEnvironment) CreateCmd(name string, args ...string) (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	cmd := exec.Command(name, args...)
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
func (r *RuntimeEnvironment) CreateWaitGroup(stdout io.ReadCloser, fOut *os.File, stderr io.ReadCloser, fErr *os.File) *sync.WaitGroup {
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
		writer := io.MultiWriter(fOut, os.Stderr)
		for scanner.Scan() {
			fmt.Fprintln(writer, scanner.Text())
		}
		wg.Done()
	}()
	return wg
}
