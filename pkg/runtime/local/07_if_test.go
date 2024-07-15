package localruntime

import (
	"bytes"
	"charlotte/pkg/job"
	"log"
	"os"
	"testing"
)

func TestIf(t *testing.T) {
	b, err := os.ReadFile("tests/07_if.yml")
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
  if !bytes.Equal(o1, []byte("")) {
    t.Errorf("invalid bytes on stdout for step 1")
  }
  e1, _ := jobRunResult.GetStepStderr(1)
  if !bytes.Equal(e1, []byte("")) {
    t.Errorf("invalid bytes on stderr for step 1")
  }

  e2, _ := jobRunResult.GetStepStderr(2)
  if !bytes.Equal(e2, []byte("[You should see this one]\n")) {
    t.Errorf("invalid bytes on stderr for step 2")
  }

  o3, _ := jobRunResult.GetStepStdout(3)
  if !bytes.Equal(o3, []byte("[Also, this one]\n")) {
    t.Errorf("invalid bytes on stdout for step 3")
  }

  o4, _ := jobRunResult.GetStepStdout(4)
  if !bytes.Equal(o4, []byte("[Again, this one]\n")) {
    t.Errorf("invalid bytes on stdout for step 4")
  }
}

