package main

import (
	"charlotte/pkg/job"
	localruntime "charlotte/pkg/runtime/local"
	"encoding/json"
	"fmt"
	"os"

	gocli "github.com/nicholasgasior/go-broccli"
)

func main() {
	cli := gocli.NewCLI("job", "Run job locally", "Streamln <hello@streamln.dev>")
	cmdRun := cli.AddCmd("run-local", "Runs YAML job file", runJobHandler)
	cmdRun.AddFlag("input-file", "i", "FILENAME", "Path to filename with a job", gocli.TypePathFile, gocli.IsRequired|gocli.IsExistent|gocli.IsRegularFile)
	cmdRun.AddFlag("quiet", "q", "", "Do not print step stdout and stderr", gocli.TypeBool, 0, gocli.OnTrue(func(c *gocli.Cmd) {

	}))
	cmdRun.AddFlag("output-file", "o", "FILENAME", "Path to write JSON result", gocli.TypePathFile, 0)
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

func runJobHandler(c *gocli.CLI) int {
	inputFile := c.Flag("input-file")
	j, err := job.NewFromFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing job from file %s: %s\n", inputFile, err.Error())
		return 1
	}
	quiet := false
	if c.Flag("quiet") == "true" {
		quiet = true
	}
	runenv := localruntime.NewLocalRuntime(quiet)
	jobRunResult := j.Run(runenv, nil)
	if !quiet && !jobRunResult.Success {
		fmt.Fprintf(os.Stderr, "Job %s failed locally with: %s\n", inputFile, jobRunResult.Error.Error())
	}

	resultJson, err := json.Marshal(jobRunResult)
	if err != nil {
		if !quiet {
			fmt.Fprintf(os.Stdout, "Marshalling run result to JSON failed\n")
		}
		return 4
	}

	outputFile := c.Flag("output-file")
	if outputFile != "" {
		err = os.WriteFile(outputFile, resultJson, 0600)
		if err != nil {
			if !quiet {
				fmt.Fprintf(os.Stdout, "Writing JSON run result to file failed\n")
			}
			return 5
		}
	} else {
		fmt.Fprintf(os.Stdout, "%s", resultJson)
	}

	if !jobRunResult.Success {
		return 1
	}

	return 0
}
