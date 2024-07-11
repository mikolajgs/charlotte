package localruntime

import (
	"bytes"
	"charlotte/pkg/job"
	"log"
	"testing"
)

var testPostSuccessJob = `name: Test
description: Workflow with post steps
inputs:
  some_input:
    default: Joe
steps:
  - type: shell
    name: Step 1
    script: |
      echo "Do nothing";

  - type: shell
    name: Step 2
    if: '{{ .Success }}'
    run_always: true
    script: |
      echo "You should see this";

  - type: shell
    if: '{{ not .Success }}'
    run_always: true
    name: Step 3
    script: |
      echo "You should not see this";

  - type: shell
    run_always: true
    name: Step 4
    script: |
      echo "You should see that";
`

var testPostFailureJob = `name: Test
description: Workflow with post steps
inputs:
  some_input:
    default: Joe
steps:
  - type: shell
    name: Step 1
    script: |
      exit 4;

  - type: shell
    name: Step 2
    if: '{{ .Success }}'
    run_always: true
    script: |
      echo "You should not see this";

  - type: shell
    if: '{{ not .Success }}'
    run_always: true
    name: Step 3
    script: |
      echo "You should see this";

  - type: shell
    run_always: true
    name: Step 4
    script: |
      echo "You should see that";
`

func TestPostSuccess(t *testing.T) {
	j, err := job.NewFromBytes([]byte(testPostSuccessJob))
	if err != nil {
		log.Fatal(err)
	}

	runenv := &LocalRuntime{}
	jobRunResult := j.Run(runenv, nil)
	if jobRunResult.Error != nil {
		t.Errorf("job run should not fail with error")
	}

  o1, _ := jobRunResult.GetStepStdout(1)
  if !bytes.Equal(o1, []byte("You should see this\n")) {
    t.Errorf("invalid bytes on stdout for step 1")
  }

  o2, _ := jobRunResult.GetStepStdout(2)
  if !bytes.Equal(o2, []byte("")) {
    t.Errorf("invalid bytes on stdout for step 2")
  }

  o3, _ := jobRunResult.GetStepStdout(3)
  if !bytes.Equal(o3, []byte("You should see that\n")) {
    t.Errorf("invalid bytes on stdout for step 3")
  }
}

func TestPostFailure(t *testing.T) {
	j, err := job.NewFromBytes([]byte(testPostFailureJob))
	if err != nil {
		log.Fatal(err)
	}

	runenv := &LocalRuntime{}
	jobRunResult := j.Run(runenv, nil)
	if jobRunResult.Error == nil {
		t.Errorf("job run should fail with error")
	}

  o1, _ := jobRunResult.GetStepStdout(1)
  if !bytes.Equal(o1, []byte("")) {
    t.Errorf("invalid bytes on stdout for step 1")
  }

  o2, _ := jobRunResult.GetStepStdout(2)
  if !bytes.Equal(o2, []byte("You should see this\n")) {
    t.Errorf("invalid bytes on stdout for step 2")
  }

  o3, _ := jobRunResult.GetStepStdout(3)
  if !bytes.Equal(o3, []byte("You should see that\n")) {
    t.Errorf("invalid bytes on stdout for step 3")
  }
}
