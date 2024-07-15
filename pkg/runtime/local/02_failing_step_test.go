package localruntime

import (
	"charlotte/pkg/job"
	"log"
	"os"
	"testing"
)

func TestFailingStep(t *testing.T) {
	b, err := os.ReadFile("tests/02_failing_step.yml")
	if err != nil {
		log.Fatal(err)
	}

	j, err := job.NewFromBytes(b)
	if err != nil {
		log.Fatal(err)
	}

	runenv := &LocalRuntime{}
	jobRunResult := j.Run(runenv, nil)
	if jobRunResult.Error == nil {
		t.Errorf("job run should fail with error")
	}

	// Still, there should be 3 step results
	if len(jobRunResult.StepRunResults) != 3 {
		t.Errorf("invalid number of step run results")
	}

	for i, s := range []string{
		"",
		"Stderr\n",
	} {
		eb, err := jobRunResult.GetStepStderr(i)
		if err != nil {
			t.Errorf("GetStepStderr failed to return a step %d stderr", i)
		}
		if string(eb) != s {
			t.Errorf("stderr for step %d has invalid contents: %s", i, string(eb))
		}
	}

	for i, s := range []string {
		"Do nothing in Step 1\n",
		"Step 2.1;Step 2.2;Step 2.3;Step 2.4;Step 2.5;\n", // TODO: Actually, there shouldn't be \n at the end but there is
	} {
		ob, err := jobRunResult.GetStepStdout(i)
		if err != nil {
			t.Errorf("GetStepStdout failed to return a step %d stdout", i)
		}
		if string(ob) != s {
			t.Errorf("stdout for step %d has invalid contents: %s", i, string(ob))
		}
	}

	if !jobRunResult.StepRunResults[2].Skipped {
		t.Errorf("step 2 should be marked as Skipped")
	}
}
