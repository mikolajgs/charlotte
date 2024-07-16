package main

type JobRun struct {
	ID *int64 `json:"id"`
	CreatedAt *string `json:"created_at"`
	StartedAt *string `json:"started_at"`
	FinishedAt *string `json:"finished_at"`
	Result *string `json:"result"`
	Content *string `json:"content"`
}
