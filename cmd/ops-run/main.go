package main

import (
	"charlotte/pkg/job"
	dockerruntimeenvironment "charlotte/pkg/runtime-environment/docker"
	localruntimeenvironment "charlotte/pkg/runtime-environment/local"
	"fmt"
	"os"

	gocli "github.com/mikogs/go-broccli/v2"
)

func main() {
	cli := gocli.NewCLI("ops-run", "Test run", "Streamln Co. <streamln@streamln.co>")
	cmdRun := cli.AddCmd("run-job", "Runs YAML job file", runJobHandler)
	cmdRun.AddFlag("file", "f", "FILENAME", "Path to filename with a job", gocli.TypePathRegularFile|gocli.Required|gocli.MustExist, nil)
	_ = cli.AddCmd("version", "Prints version", versionHandler)
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"App", "version"}
	}
	os.Exit(cli.Run(os.Stdout, os.Stderr))
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

	runenv := &localruntimeenvironment.LocalRuntimeEnvironment{}

	exitError, errors := j.Run(runenv)
	fmt.Fprintf(os.Stdout, "exitCode: %v\n", exitError)
	fmt.Fprintf(os.Stdout, "errors: %v\n", errors)

	runenv2 := &dockerruntimeenvironment.DockerRuntimeEnvironment{}

	exitError, errors = j.Run(runenv2)
	fmt.Fprintf(os.Stdout, "exitCode: %v\n", exitError)
	fmt.Fprintf(os.Stdout, "errors: %v\n", errors)

	return 0
}
