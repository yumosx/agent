package agent

type Executor interface {
	GetType() string
	Run(step string) error
}

type Agent interface {
	Think() (bool, error)
	Act() (string, error)
}
