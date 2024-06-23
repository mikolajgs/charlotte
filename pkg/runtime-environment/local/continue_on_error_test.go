package localruntimeenvironment

import (
	"charlotte/pkg/job"
	"log"
	"testing"
)

var testContinueOnErrorJob = `name: Test
description: Workflow with bash script steps
steps:
  - type: shellScript
    name: Step 1
    description: Simple test step
    continue_on_error: false
    script: |
      echo "Do nothing in Step 1";

  - type: shellScript
    name: Step 2
    description: This step fails but job should continue
    continue_on_error: true
    script: |
      for i in 1 2 3 4 5; do
        echo -n "Step 2.$i;"
        sleep 2;
      done
      >&2 echo "Stderr";
      exit 4;

      echo -n "Step 2.6;"

  - type: shellScript
    name: Step 3
    description: This step should run
    script: |
      >&2 echo "Step 3 stderr";
      echo "Step 3 stdout";
`

func TestContinueOnError(t *testing.T) {
	j, err := job.NewFromBytes([]byte(testContinueOnErrorJob))
	if err != nil {
		log.Fatal(err)
	}

	runenv := &LocalRuntimeEnvironment{}
	jobRunResult := j.Run(runenv)
	if jobRunResult.Error != nil {
		log.Fatal(err)
	}

	// Check step 3 stderr and stdout
	if len(jobRunResult.StepRunResults) != 3 {
		t.Errorf("invalid number of step run results")
	}

	for i, s := range []string{
		"",
		"Stderr\n",
		"Step 3 stderr\n",
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
		"Step 3 stdout\n",
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
