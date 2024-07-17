package main

import (
	"database/sql"
	"fmt"
	"time"

	jobrun "charlotte/pkg/jobrun"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikolajgs/struct2sql"
)

// new todo:

// todo
// - get job run
// - check if job exists when updating
// - custom error object when job does not exist
// - get count of currently unfinished jobs
// - add config with limit of jobs processing at the same time
// - return error when limit is reached (though not here)
// - job would have timeout and it should get automatically failed if it takes too long

const getJobRunQuery = `SELECT id, created_at, started_at, finished_at, result, content FROM job_runs WHERE id=?;`;
const updateJobStartedQuery = `UPDATE job_runs SET started_at=? WHERE id=?;`;
const updateJobFinishedQuery = `UPDATE job_runs SET finished_at=?, result=? WHERE id=?;`;

var queryInsert string

func initDatabase(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("error opening sqlite3 db file: %w", err)
	}

	dbSchema := struct2sql.CreateTable(&jobrun.JobRun{}, &struct2sql.CreateTableOpts{
		PrependColumns: "id INTEGER NOT NULL PRIMARY KEY",
		ExcludeFields: map[string]bool {
			"ID": true,
		},
	})

	if _, err := db.Exec(dbSchema); err != nil {
		return nil, fmt.Errorf("error creating schema in db: %w", err)
	}

	return db, nil
}

func insertJobRun(db *sql.DB, content string) (int64, error) {
	now := time.Now().Format("2006-01-02 15:04:05")

	if queryInsert == "" {
		queryInsert = struct2sql.Insert(&jobrun.JobRun{
			CreatedAt: &now,
		}, &struct2sql.InsertOpts{
			PrependColumns: "id,",
			PrependValues: "NULL,",
			IncludeFields: map[string]bool{
				"CreatedAt": true,
				"JobText": true,
			},
		})
	}

	result, err := db.Exec(queryInsert, &now, content)
	if err != nil {
		return 0, fmt.Errorf("error inserting job run to db: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last inserted id: %w", err)
	}
	
	return id, nil
}

func getJobRun(db *sql.DB, id int64) (*jobrun.JobRun, error) {
	var jr jobrun.JobRun
	err := db.QueryRow(getJobRunQuery, id).Scan(&jr.ID, &jr.CreatedAt, &jr.StartedAt, &jr.FinishedAt, &jr.ResultText, &jr.JobText)
	if err != nil {
		return nil, fmt.Errorf("error getting job run from db: %w", err)
	}
	return &jr, nil
}

func updateJobStarted(db *sql.DB, id int64) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec(updateJobStartedQuery, &now)
	if err != nil {
		return fmt.Errorf("error updating job start in db: %w", err)
	}
	return nil
}

func updateJobFinished(db *sql.DB, id int64, result string) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec(updateJobFinishedQuery, &now, result)
	if err != nil {
		return fmt.Errorf("error updating job finish in db: %w", err)
	}
	return nil
}
