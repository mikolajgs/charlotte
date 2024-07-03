package job

import (
	"charlotte/pkg/step"
	"fmt"
)

func (j *Job) validateSteps() error {
	inputMap := j.getInputsMap()
	stepOutputs := map[string]map[string]string{}

	// Create object injected to templates in Steps
	templateObj := struct{
		Inputs *map[string]string
		Variables *map[string]string
		StepOutputs *map[string]map[string]string
	}{
		Inputs: &inputMap,
		Variables: &j.Variables,
		StepOutputs: &stepOutputs,
	}

	// Loop through steps
	for i, st := range j.Steps.([]step.IStep) {
		s, err := j.getTemplateValue(st.GetScript(), &templateObj)
		if err != nil {
			return fmt.Errorf("error processing step '%s' script (%d): %w", st.GetName(), i, err)
		}
		st.SetRunScript(s)

		if st.GetID() == "" {
			continue
		}

		_, ok := stepOutputs[st.GetID()]
		if ok {
			return fmt.Errorf("duplicate step id '%s' found", st.GetID())
		}

		err = j.processStepOutputs(st, &templateObj, &stepOutputs)
		if err != nil {
			return fmt.Errorf("error processing step '%s' outputs: %w", st.GetName(), err)
		}
	}

	return nil
}

func (j *Job) processStepOutputs(st step.IStep, templateObj interface{}, stepOutputs *map[string]map[string]string) error {
	outputs := st.GetOutputs()
	if len(outputs) > 0 {
		for n, o := range outputs {
			os, err := j.getTemplateValue(o, templateObj)
			if err != nil {
				return fmt.Errorf("error processing step '%s' output '%s': %w", st.GetName(), n, err)
			}

			_, ok := (*stepOutputs)[st.GetID()]
			if !ok {
				(*stepOutputs)[st.GetID()] = map[string]string{}
			}
			(*stepOutputs)[st.GetID()][n] = os
		}
	}

	return nil
} 
