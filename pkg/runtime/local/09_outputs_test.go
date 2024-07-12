package localruntime

import (
	"charlotte/pkg/job"
	"log"
	"testing"
)

var testOutputsJob = `name: Test
description: Workflow with outputs
steps:
  - type: shell
    name: Step 1
    id: step_1
    script: |
      echo -n "Jane" > $OUTPUTS_DIR/output_1

  - type: shell
    name: Step 2
    id: step_2
    if: '{{ .Success }}'
    run_always: true
    script: |
      echo "Just a message";
    outputs:
      output_2: Joe
outputs:
  output_1:
    value: '[{{ .StepOutputs.step_1.output_1 }}]'
  output_2:
    value: '[[{{ .StepOutputs.step_2.output_2 }}]]'
`

func TestOutputs(t *testing.T) {
	j, err := job.NewFromBytes([]byte(testOutputsJob))
	if err != nil {
		log.Fatal(err)
	}

	runenv := &LocalRuntime{}
	jobRunResult := j.Run(runenv, nil)
	if jobRunResult.Error != nil {
		t.Errorf("job run should not fail with error")
	}

	if jobRunResult.RunOutputs["output_1"] != "[Jane]" {
		t.Errorf("job has invalid output")
	}
	if jobRunResult.RunOutputs["output_2"] != "[[Joe]]" {
		t.Errorf("job has invalid output")
	}
}
