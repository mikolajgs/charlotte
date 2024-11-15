package localruntime

import (
	"bytes"
	"charlotte/pkg/job"
	"log"
	"os"
	"testing"
)

func TestStdoutAndStderr(t *testing.T) {
	b, err := os.ReadFile("tests/01_std_files.yml")
	if err != nil {
		log.Fatal(err)
	}

	j, err := job.NewFromBytes(b)
	if err != nil {
		log.Fatal(err)
	}

	runenv := &LocalRuntime{}
	jobRunResult := j.Run(runenv, nil)
	if jobRunResult.Error != nil {
		log.Fatal(err)
	}

	stderrFile := jobRunResult.StepRunResults[0].StderrFile
	stderrBytes, err := os.ReadFile(stderrFile)
	if err != nil {
		t.Errorf("cannot read step1 stderr file %s: %s", stderrFile, err.Error())
	}
	if string(stderrBytes) != "Step1 Standard Error Message\n" {
		t.Errorf("stderr file has invalid contents: %s", string(stderrBytes))
	}

	stdoutFile := jobRunResult.StepRunResults[0].StdoutFile
	stdoutBytes, err := os.ReadFile(stdoutFile)
	if err != nil {
		t.Errorf("cannot read step1 stdout file %s: %s", stdoutFile, err.Error())
	}
	if string(stdoutBytes) != "Step1 Standard Output Message\n" {
		t.Errorf("stdout file has invalid contents: %s", string(stdoutBytes))
	}

	// Check GetStepStderr and GetStepStdout functions
	e1, err := jobRunResult.GetStepStderr(0);
	if err != nil {
		t.Errorf("GetStepStderr failed with error: %s", err.Error())
	}
	if !bytes.Equal(e1, stderrBytes) {
		t.Errorf("GetStepStderr returns invalid bytes")
	}

	o1, err := jobRunResult.GetStepStdout(0);
	if err != nil {
		t.Errorf("GetStepStdout failed with error: %s", err.Error())
	}
	if !bytes.Equal(o1, stdoutBytes) {
		t.Errorf("GetStepStdout returns invalid bytes")
	}

	_, err = jobRunResult.GetStepStderr(1);
	if err == nil {
		t.Errorf("GetStepStderr should fail when no stderr for a step")
	}

	_, err = jobRunResult.GetStepStdout(1);
	if err == nil {
		t.Errorf("GetStepStdout should fail when no stdout for a step")
	}
}
