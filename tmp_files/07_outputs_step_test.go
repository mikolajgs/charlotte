package localruntimeenvironment

import (
	"testing"
)

var testOutputsJob = `name: Test
description: Workflow with inputs
inputs:
	first_name:
		description: First name
		required: true
		default_value: 'Charlotte'
steps:
  - type: shellScript
    name: Step 1
    description: Simple test step
		template: go
    script: |
      echo "Hello, {{ .Inputs["first_name"].Value }}" > /tmp/welcome-name.txt
			cat /tmp/welcome-name.txt
		outputs:
			welcome_name: /tmp/welcome-name.txt
			welcome_surname: /tmp/welcome-surname.txt
	- type: shellScript
		name: Step 2
		description: Next step
		script: |
			echo "This should not be executed"
`

func TestOutputFileExists(t *testing.T) {
  //...
}

func TestMissingOutputFile(t *testing.T) {
	// ..
}
