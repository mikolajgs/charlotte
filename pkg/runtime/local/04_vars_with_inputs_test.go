package localruntime

import (
	"bytes"
	"charlotte/pkg/job"
	"log"
	"testing"
)

var testVariablesWithInputsJob = `name: Test
description: Workflow with variable
variables:
  VAR1: '1234'
  VAR2: '{{ index .Inputs "input1" }}'
  VAR3: '{{ index .Inputs "input1" }}-{{ index .Inputs "input2" }}'
  VAR4: '{{ if eq (index .Inputs "input1") "Jane" }}Jane{{ else }}Joe{{ end }}'
inputs:
  input1:
    default: "Jane"
  input2:
    default: "Joe"
steps:
  - type: shell
    name: Step 1
    script: |
      echo "stdout:{{ index .Variables "VAR1" }}";
      >&2 echo "stderr:{{ index .Variables "VAR2" }}";
  - type: shell
    name: Step 2
    script: |
      echo "stdout:{{ index .Variables "VAR3" }}";
      >&2 echo "stderr:{{ index .Variables "VAR1" }}-{{ index .Variables "VAR3" }}";
`

func TestVariablesWithInputs(t *testing.T) {
	j, err := job.NewFromBytes([]byte(testVariablesWithInputsJob))
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
