package localruntime

import (
	"bytes"
	"charlotte/pkg/job"
	"log"
	"testing"
)

var testStepOutputsJob = `name: Test
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
    name: Step 2
    script: |
      echo "[{{ .StepOutputs.step_1.second_output }}]"
  
  - type: shell
    name: Step 3
    script: |
      >&2 echo "[{{ .StepOutputs.step_1.first_output }}]"
`

func TestStepOutputs(t *testing.T) {
	j, err := job.NewFromBytes([]byte(testStepOutputsJob))
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
}

