package configurator

import (
	"cmp"
	assertions "github.com/stretchr/testify/assert"
	"io"
	"slices"
	"testing"
)

func TestGetEmptyRegistry(t *testing.T) {
	assert := assertions.New(t)
	registry, err := NewModuleRegistry([]Agent{})
	assert.Nil(err)
	agents := registry.getAll()
	assert.Empty(agents)
	agent := registry.get("123")
	assert.Nil(agent)
}

func TestGetRegistry(t *testing.T) {
	agents := []Agent{
		NewAgent("1", func(r io.Reader, format string) error {
			return nil
		}),
		NewAgent("2", func(r io.Reader, format string) error {
			return nil
		}),
		NewAgent("3", func(r io.Reader, format string) error {
			return nil
		}),
	}
	assert := assertions.New(t)
	registry, err := NewModuleRegistry(agents)
	assert.Nil(err)
	agent := registry.get("1")
	assert.Equal("1", agent.moduleName())
	agent = registry.get("2")
	assert.Equal("2", agent.moduleName())
	agent = registry.get("3")
	assert.Equal("3", agent.moduleName())

	agentList := registry.getAll()
	slices.SortFunc(agentList, func(a, b Agent) int {
		return cmp.Compare(a.moduleName(), b.moduleName())
	})
	assert.Equal(agents, agentList)
}

func TestSetEmptyRegistry(t *testing.T) {
	assert := assertions.New(t)
	registry, err := NewModuleRegistry([]Agent{})
	assert.Nil(err)
	expectedAgent := NewAgent("44", func(r io.Reader, format string) error { return nil })
	registry.set("44", expectedAgent)
	agent := registry.get("44")
	assert.Equal("44", agent.moduleName())
	agents := registry.getAll()
	assert.Equal([]Agent{expectedAgent}, agents)
}

func TestDeepHierarchy(t *testing.T) {
	assert := assertions.New(t)
	agent1 := NewAgent("1", func(r io.Reader, format string) error { return nil })
	agent2 := NewAgent("2", func(r io.Reader, format string) error { return nil })
	agent3 := NewAgent("3", func(r io.Reader, format string) error { return nil })
	agent1.Require(agent2)
	agent2.Require(agent3)
	registry, err := NewModuleRegistry([]Agent{agent1})
	assert.Nil(err)
	expectedAgents := []Agent{agent1, agent2, agent3}
	agents := registry.getAll()
	slices.SortFunc(agents, func(a, b Agent) int {
		return cmp.Compare(a.moduleName(), b.moduleName())
	})
	assert.Equal(expectedAgents, agents)
}

func Test2LevelHierarchy(t *testing.T) {
	assert := assertions.New(t)
	agent1 := NewAgent("1", func(r io.Reader, format string) error { return nil })
	agent2 := NewAgent("2", func(r io.Reader, format string) error { return nil })
	agent3 := NewAgent("3", func(r io.Reader, format string) error { return nil })
	agent1.Require(agent2)
	agent1.Require(agent3)
	registry, err := NewModuleRegistry([]Agent{agent1})
	assert.Nil(err)
	expectedAgents := []Agent{agent1, agent2, agent3}
	agents := registry.getAll()
	slices.SortFunc(agents, func(a, b Agent) int {
		return cmp.Compare(a.moduleName(), b.moduleName())
	})
	assert.Equal(expectedAgents, agents)
}

func TestLoop(t *testing.T) {
	assert := assertions.New(t)
	agent1 := NewAgent("1", func(r io.Reader, format string) error { return nil })
	agent2 := NewAgent("2", func(r io.Reader, format string) error { return nil })
	agent3 := NewAgent("3", func(r io.Reader, format string) error { return nil })
	agent4 := NewAgent("4", func(r io.Reader, format string) error { return nil })

	agent1.Require(agent2)
	agent2.Require(agent3)
	agent3.Require(agent4)
	agent4.Require(agent1)

	registry, err := NewModuleRegistry([]Agent{agent1})
	assert.Nil(err)
	agents := registry.getAll()
	slices.SortFunc(agents, func(a, b Agent) int {
		return cmp.Compare(a.moduleName(), b.moduleName())
	})
	expectedAgents := []Agent{agent1, agent2, agent3, agent4}
	assert.Equal(expectedAgents, agents)
}

func TestTestReplacement(t *testing.T) {
	assert := assertions.New(t)
	agent1 := NewAgent("1", func(r io.Reader, format string) error { return nil })
	agent2 := NewAgent("1", func(r io.Reader, format string) error { return nil })

	registry, err := NewModuleRegistry([]Agent{agent1, agent2})
	assert.NotNil(err)
	assert.Nil(registry)
}

func TestTestReplacementNoError(t *testing.T) {
	assert := assertions.New(t)
	agent1 := NewAgent("1", func(r io.Reader, format string) error { return nil })

	registry, err := NewModuleRegistry([]Agent{agent1, agent1})
	assert.Nil(err)
	assert.NotNil(registry)
}
