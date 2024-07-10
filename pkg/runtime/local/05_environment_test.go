package localruntime

import (
	"bytes"
	"charlotte/pkg/job"
	"log"
	"testing"
)

var testEnvJob = `name: Test
description: Workflow with variable
environment:
  ENV1: '12-{{ index .Variables "VAR2" }}-{{ index .Inputs "input2" }}-34'
  ENV2: '34-{{ index .Variables "VAR1" }}-{{ index .Inputs "input1" }}-56'
variables:
  VAR1: '1234'
  VAR2: '{{ index .Inputs "input1" }}'
inputs:
  input1:
    default: "Jane"
  input2:
    default: "Joe"
steps:
  - type: shell
    name: Step 1
    script: |
      echo "{{ if eq (index .Environment "ENV1") "x" }}x{{ else }}{{ index .Environment "ENV2" }}{{ end }}!";
      >&2 echo "[$ENV1][$ENV2]";
  - type: shell
    name: Step 2
    environment:
      STEP_ENV1: Step1
      ENV2: 'Env2 Overridden'
    script: |
      echo "[$STEP_ENV1][$ENV2]"

`

func TestEnvWithVarsWithInputs(t *testing.T) {
	j, err := job.NewFromBytes([]byte(testEnvJob))
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
