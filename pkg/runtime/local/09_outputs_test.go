package localruntime

import (
	"charlotte/pkg/job"
	"log"
	"os"
	"testing"
)

func TestOutputs(t *testing.T) {
	b, err := os.ReadFile("tests/09_outputs.yml")
	if err != nil {
		log.Fatal(err)
	}

	j, err := job.NewFromBytes(b)
	if err != nil {
		log.Fatal(err)
	}

	runenv := &LocalRuntime{}
	jobRunResult := j.Run(runenv, nil)
	if jobRunResult.Error != nil {
		t.Errorf("job run should not fail with error")
	}

	if jobRunResult.Outputs["output_1"] != "[Jane]" {
		t.Errorf("job has invalid output")
	}
	if jobRunResult.Outputs["output_2"] != "[[Joe]]" {
		t.Errorf("job has invalid output")
	}
}
