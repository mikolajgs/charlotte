package job

import (
	"bytes"
	"text/template"
)

func (j *Job) getInputsMap() map[string]string {
	inputMap := map[string]string{}
	for name, inp := range j.Inputs {
		inputMap[name] = inp.RunValue
	}
	return inputMap
}

func (j *Job) getTemplateValue(tpl string, tplObj interface{}) (string, error) {
	buf := &bytes.Buffer{}
	t := template.Must(template.New("v").Parse(tpl))
	err := t.Execute(buf, &tplObj)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
