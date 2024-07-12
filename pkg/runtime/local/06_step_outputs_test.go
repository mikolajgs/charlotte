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

  - type: shell
    id: step_4
    name: Step 4
    script: |
      echo -n "{{ .Inputs.some_input }} Smith" > $OUTPUTS_DIR/third_output

  - type: shell
    name: Step 5
    script: |
      echo "{{ .StepOutputs.step_4.third_output }}"

  - type: shell
    name: Step 6
    id: step_6
    script: |
      echo -n "Jane" > $OUTPUTS_DIR/fourth_output
    outputs:
      fourth_output: Not overridding

  - type: shell
    name: Step 7
    script: |
      echo "{{ .StepOutputs.step_6.fourth_output }}"
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

  o4, _ := jobRunResult.GetStepStdout(4)
  if !bytes.Equal(o4, []byte("Joe Smith\n")) {
    t.Errorf("invalid bytes on stdout for step 4")
  }

  o6, _ := jobRunResult.GetStepStdout(6)
  if !bytes.Equal(o6, []byte("Jane\n")) {
    t.Errorf("invalid bytes on stdout for step 6")
  }
}

