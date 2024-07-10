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
	isConfigured(time time.Time) bool
	register()
}

func NewAgent(name string, updateCallback UpdateFunc) Agent {
	agent := &agentImpl{
		name:           name,
		parents:        make(map[string]Agent),
		childrens:      make(map[string]Agent),
		updateCallback: updateCallback,
		time:           nil,
	}
	return agent
}

type agentImpl struct {
	name           string
	parents        map[string]Agent
	childrens      map[string]Agent
	updateCallback UpdateFunc
	time           *time.Time
}

func (a *agentImpl) Require(agent Agent) error {
	_ = agent.addParent(a)
	a.childrens[agent.moduleName()] = agent
	return nil
}

func (a *agentImpl) update(r io.ReadSeeker) error {
	if a.isConfigured(time.Now()) {
		return nil
	}
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

func (a *agentImpl) moduleName() string {
	return a.name
}

func (a *agentImpl) addParent(agent Agent) error {
	a.parents[agent.moduleName()] = agent
	return nil
}

func (a *agentImpl) parentExists() bool {
	return len(a.parents) > 0
}

func (a *agentImpl) isConfigured(time time.Time) bool {
	return a.time != nil
}

func (a *agentImpl) register() {
	r := getModuleRegistry()
	r.set(a.name, a)
	for _, agent := range a.childrens {
		r.set(agent.moduleName(), agent)
	}
}
