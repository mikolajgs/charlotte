package main

import (
	"charlotte/pkg/job"
	dockerruntimeenvironment "charlotte/pkg/runtime-environment/docker"
	kubernetesruntimeenvironment "charlotte/pkg/runtime-environment/kubernetes"
	localruntimeenvironment "charlotte/pkg/runtime-environment/local"
	"fmt"
	"os"

	gocli "github.com/mikogs/go-broccli/v2"
)

func main() {
	cli := gocli.NewCLI("ops-run", "Test run", "Streamln Co. <streamln@streamln.co>")
	cmdRun := cli.AddCmd("run-job", "Runs YAML job file", runJobHandler)
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

	runenv := &localruntimeenvironment.LocalRuntimeEnvironment{}
	err = j.Run(runenv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Job %s failed locally with: %s\n", c.Flag("file"), err.Error())
	} else {
		fmt.Fprintf(os.Stdout, "Job %s succeeded locally\n", c.Flag("file"))
	}

	runenv2 := &dockerruntimeenvironment.DockerRuntimeEnvironment{}
	err = j.Run(runenv2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Job %s failed in a docker with: %s\n", c.Flag("file"), err.Error())
	} else {
		fmt.Fprintf(os.Stdout, "Job %s succeeded in a docker\n", c.Flag("file"))
	}

	runenv3 := &kubernetesruntimeenvironment.KubernetesRuntimeEnvironment{}
	err = j.Run(runenv3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Job %s failed in kubernetes with: %s\n", c.Flag("file"), err.Error())
	} else {
		fmt.Fprintf(os.Stdout, "Job %s succeeded in kubernetes\n", c.Flag("file"))
	}

	return 0
}
