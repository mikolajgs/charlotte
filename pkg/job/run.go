package job

import (
	jobrun "charlotte/pkg/jobrun"
	runtime "charlotte/pkg/runtime"
	"charlotte/pkg/step"
	steprun "charlotte/pkg/steprun"
	"fmt"
)

// Run executes a Step in a specific Runtime.
func (j *Job) Run(runtime runtime.IRuntime, inputs *jobrun.JobRunInputs) (*jobrun.JobRunResult) {
	jobRunResult := &jobrun.JobRunResult{}
	jobRunResult.StepRunResults = make([]*steprun.StepRunResult, 0)
	jobRunResult.StepsWithErrors = make([]int, 0)

	if inputs == nil {
		inputs = &jobrun.JobRunInputs{}
	}

	err := j.initBeforeRun(inputs)
	if err != nil {
		jobRunResult.Error = fmt.Errorf("error init before run: %w", err)
		return jobRunResult
	}

	// Create runtime
	err = runtime.Create(j.Steps.([]step.IStep))
	if err != nil {
		jobRunResult.Error = fmt.Errorf("error creating runtimeenv: %w", err)
		jobRunResult.Success = false
		return jobRunResult
	}
	defer runtime.Destroy(j.Steps.([]step.IStep))


	// Actual loop
	inputMap := j.getInputsMap()
	stepOutputs := map[string]map[string]string{}

	// Create object injected to templates in Steps
	templateObj := struct{
		Inputs *map[string]string
		Variables *map[string]string
		Environment *map[string]string
		StepOutputs *map[string]map[string]string
	}{
		Inputs: &inputMap,
		Variables: &j.Variables,
		Environment: &j.Environment,
		StepOutputs: &stepOutputs,
	}

	for i, st := range j.Steps.([]step.IStep) {
		stepEnvironments := map[string]string{}
		err := j.processStepEnvironment(st, &templateObj, &stepEnvironments)
		if err != nil {
			jobRunResult.Error = fmt.Errorf("error processing step '%s' environment whilst running: %w", st.GetName(), err)
			return jobRunResult
		}

		// Get script to execute by processing the go template
		s, err := j.getTemplateValue(st.GetScript(), &templateObj)
		if err != nil {
			jobRunResult.Error = fmt.Errorf("error processing step '%s' script (%d) whilst running: %w", st.GetName(), i, err)
			return jobRunResult
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
			jobRunResult.Error = fmt.Errorf("running step %s failed with: %w", st.GetName(), err)
			jobRunResult.Success = false
			jobRunResult.BreakingStep = i

			return jobRunResult
		}

		err = j.processStepOutputs(st, &templateObj, &stepOutputs)
		if err != nil {
			jobRunResult.Error = fmt.Errorf("error processing step '%s' outputs whilst running: %w", st.GetName(), err)
			return jobRunResult
		}
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
