package job

import (
	jobrun "charlotte/pkg/jobrun"
	"fmt"
	"regexp"
)

func (j *Job) processInputs(inputs *jobrun.JobRunInputs) error {
	for n, i := range j.Inputs {
		if i.Required {
			val, ok := inputs.Inputs[n]
			if !ok || val == "" {
				return fmt.Errorf("required input %s is empty", n)
			}
		}

		if i.RegExp != "" {
			val, ok := inputs.Inputs[n]
			if ok && val != "" {
				match, _ := regexp.MatchString(i.RegExp, inputs.Inputs[n])
				if !match {
					return fmt.Errorf("inputs %s does not match regexp", n)
				}
			}
		}

		val, ok := inputs.Inputs[n]
		if !ok || val == "" {
			j.Inputs[n].RunValue = i.Default
		} else {
			j.Inputs[n].RunValue = inputs.Inputs[n]
		}
	}

	return nil
}
