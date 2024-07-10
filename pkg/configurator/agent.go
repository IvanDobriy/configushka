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
	childrenExists() bool
	moduleName() string
	isConfigured(time time.Time) bool
	signUp(registry Registry)
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
	now := time.Now()
	if a.isConfigured(now) {
		return nil
	}
	var leaf Agent = nil
	for _, agent := range a.childrens {
		if !agent.isConfigured(now) {
			leaf = agent
			break
		}
	}
	if leaf != nil {
		if err := leaf.update(r); err != nil {
			return err
		}
	}
	if a.isConfigured(now) {
		return nil
	}
	r.Seek(0, 0)
	if err := a.updateCallback(r); err != nil {
		return err
	}
	a.time = &now
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

func (a *agentImpl) childrenExists() bool {
	return len(a.parents) > 0
}

func (a *agentImpl) isConfigured(time time.Time) bool {
	return a.time != nil
}

func (a *agentImpl) signUp(registry Registry) {
	registry.set(a.name, a)
	for _, agent := range a.childrens {
		registry.set(agent.moduleName(), agent)
	}
}
