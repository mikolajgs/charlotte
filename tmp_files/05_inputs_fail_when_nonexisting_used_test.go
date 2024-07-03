package localruntime

import (
	"charlotte/pkg/job"
	jobrun "charlotte/pkg/jobrun"
	"log"
	"testing"
)

var testInputsFailWhenNonExistingJob = `name: Test
description: Workflow with inputs
inputs:
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

func TestInputsFailWhenNonExistingUsed(t *testing.T) {
	j, err := job.NewFromBytes([]byte(testInputsFailWhenNonExistingJob))
	if err != nil {
		log.Fatal(err)
	}

	runenv := &LocalRuntime{}
	jobRunResult := j.Run(runenv, &jobrun.JobRunInputs{
		Inputs: map[string]string{
			"nonexisting": "123",
		},
	})
	if jobRunResult.Error == nil {
		t.Errorf("job run should fail with error of non-existing inputs")
	}
}
