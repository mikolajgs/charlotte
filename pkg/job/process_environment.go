package job

import (
	"fmt"
)

func (j *Job) processEnvironment() error {
	inputMap := j.getInputsMap()

	// Create object injected to templates in Variables
	templateObj := struct{
		Inputs *map[string]string
		Variables *map[string]string
	}{
		Inputs: &inputMap,
		Variables: &j.Variables,
	}

	// Process template for each Environment
	for n, e := range j.Environment {
		s, err := j.getTemplateValue(e, &templateObj)
		if err != nil {
			return fmt.Errorf("error processing environment %s: %w", n, err)
		}
		// Write processed variable back to Job
		j.Environment[n] = s
	}

	return nil
}
