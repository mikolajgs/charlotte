package localruntime

import (
	"bytes"
	"charlotte/pkg/job"
	"log"
	"os"
	"testing"
)

func TestPostSuccess(t *testing.T) {
	b, err := os.ReadFile("tests/08_post_success.yml")
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
  if !bytes.Equal(o1, []byte("You should see this\n")) {
    t.Errorf("invalid bytes on stdout for step 1")
  }

  o2, _ := jobRunResult.GetStepStdout(2)
  if !bytes.Equal(o2, []byte("")) {
    t.Errorf("invalid bytes on stdout for step 2")
  }

  o3, _ := jobRunResult.GetStepStdout(3)
  if !bytes.Equal(o3, []byte("You should see that\n")) {
    t.Errorf("invalid bytes on stdout for step 3")
  }
}

func TestPostFailure(t *testing.T) {
	b, err := os.ReadFile("tests/08_post_failure.yml")
	if err != nil {
		log.Fatal(err)
	}

	j, err := job.NewFromBytes(b)
	if err != nil {
		log.Fatal(err)
	}

	runenv := &LocalRuntime{}
	jobRunResult := j.Run(runenv, nil)
	if jobRunResult.Error == nil {
		t.Errorf("job run should fail with error")
	}

  o1, _ := jobRunResult.GetStepStdout(1)
  if !bytes.Equal(o1, []byte("")) {
    t.Errorf("invalid bytes on stdout for step 1")
  }

  o2, _ := jobRunResult.GetStepStdout(2)
  if !bytes.Equal(o2, []byte("You should see this\n")) {
    t.Errorf("invalid bytes on stdout for step 2")
  }

  o3, _ := jobRunResult.GetStepStdout(3)
  if !bytes.Equal(o3, []byte("You should see that\n")) {
    t.Errorf("invalid bytes on stdout for step 3")
  }
}
