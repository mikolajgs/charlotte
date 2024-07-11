package jobrun

import (
	steprun "charlotte/pkg/steprun"
	"fmt"
	"os"
)

type JobRunResult struct {
	Success bool
	StepRunResults []*steprun.StepRunResult
	BreakingStep int
	StepsWithErrors []int
	Error error
}

func (r *JobRunResult) GetStepStderr (i int) ([]byte, error) {
	if len(r.StepRunResults) < i+1 || r.StepRunResults[i] == nil {
		return []byte{}, fmt.Errorf("step run result %d not found", i)
	}

	f := r.StepRunResults[i].StderrFile
	b, err := os.ReadFile(f)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot read stderr file %s: %w", f, err)
	}

	return b, nil
}

func (r *JobRunResult) GetStepStdout (i int) ([]byte, error) {
	if len(r.StepRunResults) < i+1 || r.StepRunResults[i] == nil {
		return []byte{}, fmt.Errorf("step run result %d not found", i)
	}

	f := r.StepRunResults[i].StdoutFile
	b, err := os.ReadFile(f)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot read stdout file %s: %w", f, err)
	}

	return b, nil
}

func (r *JobRunResult) SetFailure(e error, i int) {
	r.Error = e
	r.Success = false
	r.BreakingStep = i
}

func NewJobRunResult() *JobRunResult {
	r := &JobRunResult{}
	r.StepRunResults = make([]*steprun.StepRunResult, 0)
	r.StepsWithErrors = make([]int, 0)
	r.Success = true
	return r
}
