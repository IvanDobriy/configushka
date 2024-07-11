package configurator

import (
	"errors"
	"fmt"
	"io"
	"time"
)

type UpdateFunc func(r io.Reader, format string) error

type Agent interface {
	Require(agent Agent) error
	update(r io.ReadSeeker, format string) error
	addParent(agent Agent) error
	childrenExists() bool
	moduleName() string
	isConfigured(time time.Time) bool
	signUp(registry Registry) error
}

func NewAgent(name string, updateCallback UpdateFunc) Agent {
	agent := &agentImpl{
		name:           name,
		parents:        make(map[string]Agent),
		childrens:      make(map[string]Agent),
		updateCallback: updateCallback,
		time:           nil,
		isHandled:      false,
	}
	return agent
}

type agentImpl struct {
	name           string
	parents        map[string]Agent
	childrens      map[string]Agent
	updateCallback UpdateFunc
	time           *time.Time
	isHandled      bool
}

func (a *agentImpl) Require(agent Agent) error {
	_ = agent.addParent(a)
	a.childrens[agent.moduleName()] = agent
	return nil
}

func (a *agentImpl) update(r io.ReadSeeker, format string) error {
	now := time.Now()
	if a.isHandled {
		return nil
	}
	a.isHandled = true

	var leaf Agent = nil
	for _, agent := range a.childrens {
		if !agent.isConfigured(now) {
			leaf = agent
			break
		}
	}
	if leaf != nil {
		if err := leaf.update(r, format); err != nil {
			return err
		}
	}
	if a.isConfigured(now) {
		return nil
	}

	if _, err := r.Seek(0, 0); err != nil {
		return err
	}

	if err := a.updateCallback(r, format); err != nil {
		return err
	}
	a.time = &now
	for _, agent := range a.parents {
		if err := agent.update(r, format); err != nil {
			return err
		}
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
	return len(a.childrens) > 0
}

func (a *agentImpl) isConfigured(time time.Time) bool {
	return a.time != nil
}

func (a *agentImpl) signUp(registry Registry) error {
	oldAgent := registry.get(a.name)
	if oldAgent != nil {
		if oldAgent != a {
			return errors.New(fmt.Sprintf("found different agents with same name %v", a.name))
		}
		return nil
	}
	registry.set(a.name, a)
	for _, agent := range a.childrens {
		if err := agent.signUp(registry); err != nil {
			return err
		}
	}
	return nil
}
