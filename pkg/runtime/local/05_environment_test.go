package localruntime

import (
	"bytes"
	"charlotte/pkg/job"
	"log"
	"os"
	"testing"
)

func TestEnvWithVarsWithInputs(t *testing.T) {
	b, err := os.ReadFile("tests/05_environment.yml")
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

  if j.Environment["ENV1"] != "12-Jane-Joe-34" {
    t.Errorf("environment has invalid value")
  }
  if j.Environment["ENV2"] != "34-1234-Jane-56" {
    t.Errorf("environment has invalid value")
  }

  o0, _ := jobRunResult.GetStepStdout(0)
  if !bytes.Equal(o0, []byte("34-1234-Jane-56!\n")) {
    t.Errorf("invalid bytes on stdout for step 0")
  }
  e0, _ := jobRunResult.GetStepStderr(0)
  if !bytes.Equal(e0, []byte("[12-Jane-Joe-34][34-1234-Jane-56]\n")) {
    t.Errorf("invalid bytes on stderr for step 0")
  }

  o1, _ := jobRunResult.GetStepStdout(1)
  if !bytes.Equal(o1, []byte("[Step1][Env2 Overridden]\n")) {
    t.Errorf("invalid bytes on stdout for step 1")
  }
}
