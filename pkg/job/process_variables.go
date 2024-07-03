package job

import (
	"fmt"
)

func (j *Job) processVariables() error {
	inputMap := j.getInputsMap()

	// Create object injected to templates in Variables
	templateObj := struct{
		Inputs *map[string]string
	}{
		Inputs: &inputMap,
	}

	// Process template for each Variable
	for n, v := range j.Variables {
		s, err := j.getTemplateValue(v, &templateObj)
		if err != nil {
			return fmt.Errorf("error processing variable %s: %w", n, err)
		}
		// Write processed variable back to Job
		j.Variables[n] = s
	}

	return nil
}
