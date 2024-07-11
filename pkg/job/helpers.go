package job

import (
	"bytes"
	"charlotte/pkg/step"
	"fmt"
	"text/template"
)

type TemplateObj struct {
	Inputs *map[string]string
	Variables *map[string]string
	Environment *map[string]string
	StepOutputs *map[string]map[string]string
	Success bool
}

func (j *Job) getInputsMap() map[string]string {
	inputMap := map[string]string{}
	for name, inp := range j.Inputs {
		inputMap[name] = inp.RunValue
	}
	return inputMap
}

func (j *Job) getTemplateValue(tpl string, tplObj *TemplateObj) (string, error) {
	buf := &bytes.Buffer{}
	t := template.Must(template.New("v").Parse(tpl))
	err := t.Execute(buf, &tplObj)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (j *Job) processStepOutputs(st step.IStep, templateObj *TemplateObj, stepOutputs *map[string]map[string]string) error {
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

func (j *Job) processStepEnvironment(st step.IStep, templateObj *TemplateObj, stepEnvironments *map[string]string) error {
	envVars := st.GetEnvironment()
	if len(envVars) > 0 {
		for n, ev := range envVars {
			es, err := j.getTemplateValue(ev, templateObj)
			if err != nil {
				return fmt.Errorf("error processing step '%s' output '%s': %w", st.GetName(), n, err)
			}

			(*stepEnvironments)[n] = es
		}
	}

	return nil
} 
