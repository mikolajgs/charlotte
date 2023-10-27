package step

type IStep interface {
	GetName() string
	GetDescription() string
	GetContinueOnError() bool
	Validate() (bool, map[string]int)
	GetScript() string
}
