package step

type IStep interface {
	GetName() string
	GetDescription() string
	GetContinueOnError() bool
	Validate() (bool, map[string]int)
	GetScript() string
	GetRunScript() string
	SetRunScript(string)
	GetOutputs() map[string]string
	SetOutput(string, string)
	GetID() string
	GetEnvironment() map[string]string
	SetEnvironmentVar(string, string)
	GetIf() string
}
