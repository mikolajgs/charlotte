package job

import (
	"charlotte/pkg/step"
	"fmt"
	"strings"
)

func (j *Job) validateSteps() error {
	inputMap := j.getInputsMap()
	stepOutputs := map[string]map[string]string{}

	// Create object injected to templates in Steps
	templateObj := &TemplateObj{
		Inputs: &inputMap,
		Variables: &j.Variables,
		Environment: &j.Environment,
		StepOutputs: &stepOutputs,
	}

	// Loop through steps
	for i, st := range j.Steps.([]step.IStep) {
		stepEnvironments := map[string]string{}
		err := j.processStepEnvironment(st, templateObj, &stepEnvironments)
		if err != nil {
			return fmt.Errorf("error processing step '%s' environment: %w", st.GetName(), err)
		}

		ifTpl := strings.TrimSpace(st.GetIf())
		if ifTpl != "" {
			_, err = j.getTemplateValue(ifTpl, templateObj)
			if err != nil {
				return fmt.Errorf("error processing step '%s' if (%d): %w", st.GetName(), i, err)
			}
		}

		s, err := j.getTemplateValue(st.GetScript(), templateObj)
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

		err = j.processStepOutputs(st, templateObj, &stepOutputs)
		if err != nil {
			return fmt.Errorf("error processing step '%s' outputs: %w", st.GetName(), err)
		}
	}

	return nil
}
