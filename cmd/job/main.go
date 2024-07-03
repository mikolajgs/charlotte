package main

import (
	"charlotte/pkg/job"
	localruntime "charlotte/pkg/runtime/local"
	"fmt"
	"os"

	gocli "github.com/nicholasgasior/go-broccli"
)

func main() {
	cli := gocli.NewCLI("job", "Run job locally", "Streamln <hello@streamln.dev>")
	cmdRun := cli.AddCmd("run-local", "Runs YAML job file", runJobHandler)
	cmdRun.AddFlag("file", "f", "FILENAME", "Path to filename with a job", gocli.TypePathFile, gocli.IsRequired|gocli.IsExistent|gocli.IsRegularFile)
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
	j, err := job.NewFromFile(c.Flag("file"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing job from file %s: %s\n", c.Flag("file"), err.Error())
		return 1
	}

	runenv := &localruntime.LocalRuntime{}
	jobRunResult := j.Run(runenv, nil)
	if jobRunResult.Error != nil {
		fmt.Fprintf(os.Stderr, "Job %s failed locally with: %s\n", c.Flag("file"), jobRunResult.Error.Error())
	} else {
		fmt.Fprintf(os.Stdout, "Job %s succeeded locally\n", c.Flag("file"))
	}

	return 0
}
