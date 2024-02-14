package runtimeenvironment

import (
	"bufio"
	"charlotte/pkg/step"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type IRuntimeEnvironment interface {
	Run(step step.IStep) (int, []string, string, string)
}

type RuntimeEnvironment struct {
}

func (r *RuntimeEnvironment) InitRunStep(step step.IStep) (int, string, *os.File, *os.File, *os.File) {
	script := step.GetScript()
	if script == "" {
		return 101, "script is empty", nil, nil, nil
	}

	fStep, err := os.CreateTemp("", "step.*.sh")
	if err != nil {
		return 103, fmt.Sprint("error creating tmp file for script: %s", err.Error()), nil, nil, nil
	}

	if _, err := fStep.Write([]byte(script)); err != nil {
		fStep.Close()
		return 104, fmt.Sprint("error writing tmp file: %s", err.Error()), nil, nil, nil
	}

	fOut, err := os.CreateTemp("", "stepout.*.txt")
	if err != nil {
		return 111, fmt.Sprint("error creating tmp file for stdout: %s", err.Error()), nil, nil, nil
	}

	fErr, err := os.CreateTemp("", "steperr.*.txt")
	if err != nil {
		return 112, fmt.Sprint("error creating tmp file for stderr: %s", err.Error()), nil, nil, nil
	}

	return 0, "", fStep, fOut, fErr
}

func (r *RuntimeEnvironment) CreateCmd(name string, args ...string) (int, string, *exec.Cmd, io.ReadCloser, io.ReadCloser) {
	cmd := exec.Command(name, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 106, fmt.Sprintf("error piping stdout: %s", err.Error()), nil, nil, nil
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return 107, fmt.Sprintf("error piping stderr: %s", err.Error()), nil, nil, nil
	}
	return 0, "", cmd, stdout, stderr
}

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
