package steprun

type StepRunResult struct {
	Success bool
	StderrFile string
	StdoutFile string
	Error error
	Skipped bool
}
