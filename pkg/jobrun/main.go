package jobrun

type JobRun struct {
	ID *int64 `json:"id"`
	CreatedAt *string `json:"created_at"`
	StartedAt *string `json:"started_at"`
	FinishedAt *string `json:"finished_at"`
	InputsText *string `json:"inputs"`
	ResultText *string `json:"result"`
	JobText *string `json:"job"`
}
