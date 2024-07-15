package job

import (
	jobrun "charlotte/pkg/jobrun"
	runtime "charlotte/pkg/runtime"
	"charlotte/pkg/step"
	steprun "charlotte/pkg/steprun"
	"fmt"
	"strings"
)

// Run executes a Step in a specific Runtime.
func (j *Job) Run(runtime runtime.IRuntime, inputs *jobrun.JobRunInputs) (*jobrun.JobRunResult) {
	jobRunResult := jobrun.NewJobRunResult()

	if inputs == nil {
		inputs = &jobrun.JobRunInputs{}
	}

	err := j.initBeforeRun(inputs)
	if err != nil {
		jobRunResult.SetFailure(fmt.Errorf("error init before run: %w", err), -1)
		return jobRunResult
	}

	// Create runtime
	err = runtime.Create(j.Steps.([]step.IStep))
	if err != nil {
		jobRunResult.SetFailure(fmt.Errorf("error creating runtime: %w", err), -1)
		return jobRunResult
	}
	defer runtime.Destroy(j.Steps.([]step.IStep))


	// Actual loop
	inputMap := j.getInputsMap()
	stepOutputs := map[string]map[string]string{}

	// Create object injected to templates in Steps
	templateObj := &TemplateObj{
		Inputs: &inputMap,
		Variables: &j.Variables,
		Environment: &j.Environment,
		StepOutputs: &stepOutputs,
		Success: true,
	}

	for i, st := range j.Steps.([]step.IStep) {
		// Logic that determines whether the step should be skipped or executed
		executeStep := true
		if !jobRunResult.Success && !st.GetRunAlways() {
			executeStep = false
		}

		// If step is not meant to be executed due to an error that happened earlier,
		// then its result is marked as Skipped
		if !executeStep {
			jobRunResult.StepRunResults = append(jobRunResult.StepRunResults, steprun.NewSkippedRunResult())
			templateObj.Success = false
			continue
		}
		
		// Get if condition
		ifTpl := strings.TrimSpace(st.GetIf())
		if ifTpl != "" {
			ifVal, err := j.getTemplateValue(ifTpl, templateObj)
			if err != nil {
				jobRunResult.SetFailure(fmt.Errorf("error processing step '%s' if (%d) whilst running: %w", st.GetName(), i, err), i)
				templateObj.Success = false
				continue
			}
			// If condition must return string '1' or 't'
			if ifVal != "true" {
				jobRunResult.StepRunResults = append(jobRunResult.StepRunResults, steprun.NewSkippedRunResult())
				continue
			}
		}

		// Get step environments - to be merged with global ones later on
		stepEnvironments := map[string]string{}
		err := j.processStepEnvironment(st, templateObj, &stepEnvironments)
		if err != nil {
			jobRunResult.SetFailure(fmt.Errorf("error processing step '%s' environment whilst running: %w", st.GetName(), err), i)
			templateObj.Success = false
			continue
		}

		// Get script to execute by processing the go template
		s, err := j.getTemplateValue(st.GetScript(), templateObj)
		if err != nil {
			jobRunResult.SetFailure(fmt.Errorf("error processing step '%s' script (%d) whilst running: %w", st.GetName(), i, err), i)
			templateObj.Success = false
			continue
		}
		st.SetRunScript(s)

		// Merge global and step environment variables
		runEnv := map[string]string{}
		for k, globalVal := range j.Environment {
			runEnv[k] = globalVal
		}
		for k, stepVal := range stepEnvironments {
			runEnv[k] = stepVal
		}
	
		// Execute step
		fOut, fErr, err := runtime.Run(st, i, &runEnv)

		suc := true
		if err != nil {
			suc = false
			jobRunResult.StepsWithErrors = append(jobRunResult.StepsWithErrors, i)
		}

		jobRunResult.StepRunResults = append(jobRunResult.StepRunResults, &steprun.StepRunResult{
			Success: suc,
			StderrFile: fErr,
			StdoutFile: fOut,
			Error: err,
		})

		if err != nil && !st.GetContinueOnError() {
			jobRunResult.SetFailure(fmt.Errorf("running step %s failed with: %w", st.GetName(), err), i)
			templateObj.Success = false
			continue
		}

		err = j.processStepOutputs(st, templateObj, &stepOutputs)
		if err != nil {
			jobRunResult.SetFailure(fmt.Errorf("error processing step '%s' outputs whilst running: %w", st.GetName(), err), i)
			templateObj.Success = false
			continue
		}
	}

	// Process job outputs if success
	jobRunResult.Outputs = map[string]string{}
	if jobRunResult.Success {
		for n, out := range j.Outputs {
			s, err := j.getTemplateValue(out.Value, templateObj)
			if err != nil {
				jobRunResult.SetFailure(fmt.Errorf("error processing output '%s': %w", n, err), 0)
				continue
			}
			jobRunResult.Outputs[n] = s
		}
	}

	if jobRunResult.Error != nil {
		jobRunResult.ErrorString = jobRunResult.Error.Error()
	}

	return jobRunResult
}

func (j *Job) initBeforeRun(inputs *jobrun.JobRunInputs) error {
	// Validate inputs and assigned their values to Job object
	err := j.processInputs(inputs)
	if err != nil {
		return err
	}

	// Process variables that may contain inputs
	err = j.processVariables()
	if err != nil {
		return err
	}

	// Process environment variables that might use inputs and variables
	err = j.processEnvironment()
	if err != nil {
		return err
	}

	// Pre-process (but do not run) steps to validate their contents
	err = j.validateSteps()
	if err != nil {
		return err
	}

	return nil
}
