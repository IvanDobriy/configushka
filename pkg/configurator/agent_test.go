package configurator

import (
	assertions "github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
	"time"
)

func TestModuleName(t *testing.T) {
	assert := assertions.New(t)
	agent := NewAgent("1", func(r io.Reader, format string) error {
		return nil
	})
	assert.Equal("1", agent.moduleName())
}

func TestAddParent(t *testing.T) {
	assert := assertions.New(t)
	agent := &agentImpl{
		name:      "1",
		parents:   make(map[string]Agent),
		childrens: make(map[string]Agent),
		updateCallback: func(r io.Reader, format string) error {
			return nil
		},
		time: nil,
	}
	agent1 := NewAgent("2", func(r io.Reader, format string) error {
		return nil
	})
	_ = agent.addParent(agent1)
	assert.Equal(agent1, agent.parents["2"])
}

func TestRequire(t *testing.T) {
	assert := assertions.New(t)
	agent1 := &agentImpl{
		name:      "1",
		parents:   make(map[string]Agent),
		childrens: make(map[string]Agent),
		updateCallback: func(r io.Reader, format string) error {
			return nil
		},
		time: nil,
	}
	agent2 := &agentImpl{
		name:      "2",
		parents:   make(map[string]Agent),
		childrens: make(map[string]Agent),
		updateCallback: func(r io.Reader, format string) error {
			return nil
		},
		time: nil,
	}
	_ = agent1.Require(agent2)
	assert.Equal(agent2, agent1.childrens["2"])
	assert.Equal(agent1, agent2.parents["1"])
}

func TestChildrenExists(t *testing.T) {
	assert := assertions.New(t)
	agent1 := NewAgent("1", func(r io.Reader, format string) error {
		return nil
	})
	agent2 := NewAgent("2", func(r io.Reader, format string) error {
		return nil
	})
	_ = agent1.Require(agent2)
	assert.True(agent1.childrenExists())
}

func TestIsConfigured(t *testing.T) {
	assert := assertions.New(t)
	agent1 := NewAgent("1", func(r io.Reader, format string) error { return nil })
	agent2 := NewAgent("2", func(r io.Reader, format string) error { return nil })
	_ = agent1.Require(agent2)
	now := time.Now()
	assert.False(agent1.isConfigured(now))
	assert.False(agent2.isConfigured(now))

	_ = agent1.Require(agent2)
	buffer := strings.NewReader("hello, world")
	reader := io.NewSectionReader(buffer, 0, buffer.Size())
	agent2.update(reader, "123")
	assert.True(agent1.isConfigured(now))
	assert.True(agent2.isConfigured(now))
}
