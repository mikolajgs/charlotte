package main

import (
	"charlotte/pkg/job"
	jobrun "charlotte/pkg/jobrun"
	localruntime "charlotte/pkg/runtime/local"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mikolajgs/broccli"
)
     
func main() {
	cli := broccli.NewCLI("job", "Run job locally", "Streamln <hello@streamln.dev>")
	cmdRun := cli.AddCmd("run-local", "Runs YAML job file", runJobHandler)
	cmdRun.AddFlag("job", "j", "FILENAME", "Path to filename with a job", broccli.TypePathFile, broccli.IsRequired|broccli.IsExistent|broccli.IsRegularFile)
	cmdRun.AddFlag("inputs", "i", "FILENAME", "Path to file containing input values", broccli.TypePathFile, broccli.IsRequired|broccli.IsExistent|broccli.IsRegularFile)
	cmdRun.AddFlag("quiet", "q", "", "Do not print step stdout and stderr", broccli.TypeBool, 0, broccli.OnTrue(func(c *broccli.Cmd) {

	}))
	cmdRun.AddFlag("result", "r", "FILENAME", "Path to write JSON result", broccli.TypePathFile, 0)
	_ = cli.AddCmd("version", "Prints version", versionHandler)
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"App", "version"}
	}
	os.Exit(cli.Run())
}

func versionHandler(c *broccli.CLI) int {
	fmt.Fprintf(os.Stdout, VERSION+"\n")
	return 0
}

func runJobHandler(c *broccli.CLI) int {
	jobFile := c.Flag("job")
	j, err := job.NewFromFile(jobFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing job from file %s: %s\n", jobFile, err.Error())
		return 1
	}

	inputsFile := c.Flag("inputs")
	b, err := os.ReadFile(inputsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading inputs file %s: %s\n", inputsFile, err.Error())
		return 1
	}
	var jobRunInputs jobrun.JobRunInputs
	err = json.Unmarshal(b, &jobRunInputs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshalling inputs file %s: %s\n", inputsFile, err.Error())
		return 1
	}

	quiet := false
	if c.Flag("quiet") == "true" {
		quiet = true
	}
	runenv := localruntime.NewLocalRuntime(quiet)
	jobRunResult := j.Run(runenv, &jobRunInputs)
	if !quiet && !jobRunResult.Success {
		fmt.Fprintf(os.Stderr, "Job %s failed locally with: %s\n", jobFile, jobRunResult.Error.Error())
	}

	resultJson, err := json.Marshal(jobRunResult)
	if err != nil {
		if !quiet {
			fmt.Fprintf(os.Stdout, "Marshalling run result to JSON failed\n")
		}
		return 4
	}

	outputFile := c.Flag("result")
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
