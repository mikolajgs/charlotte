package localruntime

import (
	"bytes"
	"charlotte/pkg/job"
	"log"
	"testing"
)

var testIfJob = `name: Test
description: Workflow with step outputs
inputs:
  some_input:
    default: Joe
steps:
  - type: shell
    name: Step 1
    id: step_1
    script: 'echo "Do nada!";'
    outputs:
      first_output: '123'
      second_output: 'name is {{ .Inputs.some_input }}'

  - type: shell
    if: '{{ eq .Inputs.some_input "Jane" }}'
    name: Step 2
    script: |
      echo "!!! If you can see this then something went wrong"
  
  - type: shell
    if: '{{ and (eq .Inputs.some_input "Joe") (eq .StepOutputs.step_1.first_output "123") }}'
    name: Step 3
    script: |
      >&2 echo "[You should see this one]"

  - type: shell
    name: Step 4
    script: |
      echo "[Also, this one]"

  - type: shell
    if: 'true'
    name: Step 5
    script: |
      echo "[Again, this one]"
`

func TestIf(t *testing.T) {
	j, err := job.NewFromBytes([]byte(testIfJob))
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

