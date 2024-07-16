package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	gocli "github.com/nicholasgasior/go-broccli"
)

func main() {
	cli := gocli.NewCLI("agent", "Receive jobs and execute them", "Streamln <hello@streamln.dev>")
	cmdStart := cli.AddCmd("start", "Starts daemon", startHandler)
	cmdStart.AddFlag("http-port", "p", "", "Port for HTTP daemon", gocli.TypeInt, gocli.IsRequired)
	cmdStart.AddFlag("db-file", "d", "", "Path to database file", gocli.TypeString, gocli.IsRequired)
	_ = cli.AddCmd("version", "Prints version", versionHandler)
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"App", "version"}
	}
	os.Exit(cli.Run())
}

func versionHandler(c *gocli.CLI) int {
	fmt.Fprintf(os.Stdout, VERSION+"\n")
	return 0
}

func startHandler(c *gocli.CLI) int {
	db, err := initDatabase(c.Flag("db-file"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initializing database: %s\n", err.Error())
		return -1
	}

	done := make(chan bool)
	go func() {
		http.HandleFunc("/jobruns/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				id, err := insertJobRun(db, string(body));
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(fmt.Sprintf(`{"id":%d}`, id)))
				return
			}

			if r.Method == http.MethodGet {
				id := getIDFromURI(r.RequestURI)
				if id == "" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				idInt, err := strconv.Atoi(id)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				jr, err := getJobRun(db, int64(idInt))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				jsonJobResult, err := json.Marshal(&jr)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				
				w.WriteHeader(http.StatusOK)
				w.Write(jsonJobResult)
				return
			}
		})
		if err := http.ListenAndServe(fmt.Sprintf(":%s", c.Flag("http-port")), nil); err != nil {
			panic(err)
		}
	}()

	<-done
	return 0
}

func getIDFromURI(uri string) string {
	xs := strings.SplitN(uri, "?", 2)
	if xs[0] == "" {
		return ""
	}
	id := strings.Replace(xs[0], "/jobruns/", "", 1)
	match, err := regexp.MatchString("^[0-9]+$", id)
	if err != nil {
		return ""
	}
	if !match {
		return ""
	}
	return id
}
