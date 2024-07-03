package localruntimeenvironment

var testOutputsBetweenStepsJob = `name: Test
description: Workflow with inputs
output:
	welcome_both:
		value: "{{ .Steps["step3"].Outputs["welcome_both"].Value }}"
		template: go
inputs:
	first_name:
		description: First name
		required: true
		default_value: 'Charlotte'
steps:
  - type: shellScript
    name: Step 1
		id: step1
    description: Simple test step
		template: go
    script: |
      echo "Hello, {{ .Inputs["first_name"].Value }}" > /tmp/welcome-name.txt
			cat /tmp/welcome-name.txt
		outputs:
			welcome_name: /tmp/welcome-name.txt
	- type: shellScript
		name: Step 2
		description: Next step
		script: |
			echo "This step does nothing"
	- type: shellScript
		name: Step 3
		description: Returns another output
		script: |
			echo "Hello both, {{ .Steps["step1"].Outputs["welcome_name"].Value }} and {{ .Inputs["second_name"].Value }}" > /tmp/welcome-both.txt
			cat /tmp/welcome-both.txt
		outputs:
			welcome_both: /tmp/welcome-both.txt
`
