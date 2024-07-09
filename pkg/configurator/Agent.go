package configurator

import "io"

type Agent interface {
	Require(agent Agent) error
	Update(r io.Reader) error
}

type AgentImpl struct {
}

func (a *AgentImpl) Require(agent Agent) error {
	return nil
}
func (a *AgentImpl) Update(r io.Reader) error {
	return nil
}
