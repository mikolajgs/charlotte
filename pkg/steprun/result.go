package steprun

type StepRunResult struct {
	Success bool
	StderrFile string
	StdoutFile string
	Error error
	Skipped bool
}

func NewSkippedRunResult() *StepRunResult {
	return &StepRunResult{
		Success: false,
		StderrFile: "",
		StdoutFile: "",
		Error: nil,
		Skipped: true,
	}
}
