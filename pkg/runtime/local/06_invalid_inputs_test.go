package localruntime

/*
import (
	"charlotte/pkg/job"
	"log"
	"testing"
)

var testInvalidInputInEnvJob = `name: Test
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
    script: 'echo "Do nada!";'
`

func TestInvalidInputs(t *testing.T) {
	j, err := job.NewFromBytes([]byte(testInvalidInputsJob))
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
}
*/
