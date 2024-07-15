package localruntime

import (
	"charlotte/pkg/job"
	jobrun "charlotte/pkg/jobrun"
	"log"
	"testing"
)

var testInputsGoTplJob = `name: Test
description: Workflow with inputs
inputs:
	first_name:
		description: First name
		required: true
		default_value: 'Charlotte'
steps:
  - type: shellScript
    name: Step 1
    description: Simple test step
		template: go
    script: |
      echo "Hello, {{ .Inputs["first_name"].Value }}";
`

func TestInputsGoTemplateReplace(t *testing.T) {
	j, err := job.NewFromBytes([]byte(testInputsGoTplJob))
	if err != nil {
		log.Fatal(err)
	}

	runenv := &LocalRuntime{}
	jobRunResult := j.Run(runenv, &jobrun.JobRunInputs{
		Inputs: map[string]string{
			"first_name": "Jane",
		},
	})
	if jobRunResult.Error != nil {
		t.Errorf("job run should not fail")
	}

	o1, err := jobRunResult.GetStepStdout(0);
	if err != nil {
		t.Errorf("GetStepStdout failed with error: %s", err.Error())
	}
	if string(o1) != "Hello, Jane\n" {
		t.Errorf("GetStepStdout returns invalid bytes")
	}
}

func TestInputsGoTemplateReplaceWithDefault(t *testing.T) {
	j, err := job.NewFromBytes([]byte(testInputsGoTplJob))
	if err != nil {
		log.Fatal(err)
	}

	runenv := &LocalRuntime{}
	jobRunResult := j.Run(runenv, nil)
	if jobRunResult.Error != nil {
		t.Errorf("job run should not fail")
	}

	o1, err := jobRunResult.GetStepStdout(0);
	if err != nil {
		t.Errorf("GetStepStdout failed with error: %s", err.Error())
	}
	if string(o1) != "Hello, Charlotte\n" {
		t.Errorf("GetStepStdout returns invalid bytes")
	}
}
