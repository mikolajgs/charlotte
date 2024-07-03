package localruntime

import (
	"charlotte/pkg/job"
	"log"
	"testing"
)

var testFailingJob = `name: Test
description: Workflow with bash script steps
steps:
  - type: shell
    name: Step 1
    description: Simple test step
    script: |
      echo "Do nothing in Step 1";

  - type: shell
    name: Step 2
    description: This step fails and it should not continue further
    script: |
      for i in 1 2 3 4 5; do
        echo -n "Step 2.$i;"
        sleep 2;
      done
      >&2 echo "Stderr";
      exit 4;

      echo -n "Step 2.6;"

  - type: shell
    name: Step 3
    description: This step should run
    script: |
      >&2 echo "Step 3 stderr";
      echo "Step 3 stdout";
`

func TestFailingStep(t *testing.T) {
	j, err := job.NewFromBytes([]byte(testFailingJob))
	if err != nil {
		log.Fatal(err)
	}

	runenv := &LocalRuntime{}
	jobRunResult := j.Run(runenv, nil)
	if jobRunResult.Error == nil {
		t.Errorf("job run should fail with error")
	}

	// There should be only 2 steps
	if len(jobRunResult.StepRunResults) != 2 {
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
}
