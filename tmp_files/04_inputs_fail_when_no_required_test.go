package localruntime

import (
	"charlotte/pkg/job"
	jobrun "charlotte/pkg/jobrun"
	"log"
	"testing"
)

var testInputsFailWhenMissingJob = `name: Test
description: Workflow with inputs
inputs:
	numeric_value:
		description: Numeric input
		required: false
		default_value: '4'
	string_value:
		description: String input
		required: true
		default_value: 'charlotte'
steps:
  - type: shell
    name: Step 1
    description: Simple test step
    script: |
      echo "Do nothing in Step 1";
`

func TestInputsFailWhenRequiredMissing(t *testing.T) {
	j, err := job.NewFromBytes([]byte(testInputsFailWhenMissingJob))
	if err != nil {
		log.Fatal(err)
	}

	runenv := &LocalRuntime{}
	jobRunResult := j.Run(runenv, &jobrun.JobRunInputs{
		Inputs: map[string]string{
			"numeric_value": "123",
		},
	})
	if jobRunResult.Error == nil {
		t.Errorf("job run should fail with error of missing inputs")
	}
}
