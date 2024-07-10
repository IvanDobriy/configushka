package configurator

import (
	"io"
	"time"
)

type UpdateFunc func(r io.Reader) error

type Agent interface {
	Require(agent Agent) error
	update(r io.ReadSeeker) error
	addParent(agent Agent) error
	parentExists() bool
	moduleName() string
}

func NewAgent(name string, updateCallback UpdateFunc) Agent {
	agent := &AgentImpl{
		name:           name,
		parents:        make(map[string]Agent),
		childrens:      make(map[string]Agent),
		updateCallback: updateCallback,
		time:           nil,
	}
	r := NewRegistry()
	r.Set(agent.name, agent)
	return agent
}

type AgentImpl struct {
	name           string
	parents        map[string]Agent
	childrens      map[string]Agent
	updateCallback UpdateFunc
	time           *time.Time
}

func (a *AgentImpl) Require(agent Agent) error {
	_ = agent.addParent(a)
	a.childrens[agent.moduleName()] = agent
	return nil
}

func (a *AgentImpl) update(r io.ReadSeeker) error {
	var leaf Agent = nil
	for _, agent := range a.childrens {
		if agent.parentExists() {
			leaf = agent
			break
		}
	}
	if leaf == nil {
		return nil
	}
	r.Seek(0, 0)
	if err := a.updateCallback(r); err != nil {
		return err
	}
	for _, agent := range a.parents {
		r.Seek(0, 0)
		agent.update(r)
	}
	return nil
}

func (a *AgentImpl) moduleName() string {
	return a.name
}

func (a *AgentImpl) addParent(agent Agent) error {
	a.parents[agent.moduleName()] = agent
	return nil
}

func (a *AgentImpl) parentExists() bool {
	return len(a.parents) > 0
}
