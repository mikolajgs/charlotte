package localruntime

import (
	"bytes"
	"charlotte/pkg/job"
	"log"
	"os"
	"testing"
)

func TestStepOutputs(t *testing.T) {
	b, err := os.ReadFile("tests/06_step_outputs.yml")
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

  o1, _ := jobRunResult.GetStepStdout(1)
  if !bytes.Equal(o1, []byte("[name is Joe]\n")) {
    t.Errorf("invalid bytes on stdout for step 0")
  }
  e2, _ := jobRunResult.GetStepStderr(2)
  if !bytes.Equal(e2, []byte("[123]\n")) {
    t.Errorf("invalid bytes on stderr for step 1")
  }

  o4, _ := jobRunResult.GetStepStdout(4)
  if !bytes.Equal(o4, []byte("Joe Smith\n")) {
    t.Errorf("invalid bytes on stdout for step 4")
  }

  o6, _ := jobRunResult.GetStepStdout(6)
  if !bytes.Equal(o6, []byte("Jane\n")) {
    t.Errorf("invalid bytes on stdout for step 6")
  }
}

