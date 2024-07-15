package localruntime

import (
	"bytes"
	"charlotte/pkg/job"
	"log"
	"os"
	"testing"
)

func TestVariablesWithInputs(t *testing.T) {
	b, err := os.ReadFile("tests/04_variables.yml")
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
		t.Errorf("job run should not fail with error")
	}

  if j.Variables["VAR1"] != "1234" {
    t.Errorf("variable has invalid value")
  }
  if j.Variables["VAR2"] != "Jane" {
    t.Errorf("variable has invalid value")
  }
  if j.Variables["VAR3"] != "Jane-Joe" {
    t.Errorf("variable has invalid value")
  }
  if j.Variables["VAR4"] != "Jane" {
    t.Errorf("variable has invalid value")
  }

  o0, _ := jobRunResult.GetStepStdout(0)
  if !bytes.Equal(o0, []byte("stdout:1234\n")) {
    t.Errorf("invalid bytes on stdout for step 0")
  }
  e0, _ := jobRunResult.GetStepStderr(0)
  if !bytes.Equal(e0, []byte("stderr:Jane\n")) {
    t.Errorf("invalid bytes on stderr for step 0")
  }

  o1, _ := jobRunResult.GetStepStdout(1)
  if !bytes.Equal(o1, []byte("stdout:Jane-Joe\n")) {
    t.Errorf("invalid bytes on stdout for step 1")
  }
  e1, _ := jobRunResult.GetStepStderr(1)
  if !bytes.Equal(e1, []byte("stderr:1234-Jane-Joe\n")) {
    t.Errorf("invalid bytes on stderr for step 1")
  }
}
