package agent

type PrintExecutor struct {
	maxStep int
}

func NewPrintExecutor() *PrintExecutor {
	return &PrintExecutor{}
}

func (p *PrintExecutor) GetType() string {
	return "printer"
}

func (p *PrintExecutor) Run(step string) error {
	s, err := p.Step(step)
	if err != nil {
		return err
	}
	_ = s
	return nil
}

func (p *PrintExecutor) Step(step string) (string, error) {
	println(step)
	return "", nil
}

func (p *PrintExecutor) Think() (bool, error) {
	return true, nil
}

func (p *PrintExecutor) Act() (string, error) {
	return "", nil
}

func (p *PrintExecutor) newTool() {

}
