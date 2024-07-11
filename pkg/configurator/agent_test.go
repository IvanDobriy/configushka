package configurator

import (
	"bufio"
	"cmp"
	"errors"
	assertions "github.com/stretchr/testify/assert"
	"io"
	"slices"
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

func TestSignUp(t *testing.T) {
	assert := assertions.New(t)
	agent1 := NewAgent("1", func(r io.Reader, format string) error { return nil })
	agent2 := NewAgent("2", func(r io.Reader, format string) error { return nil })
	_ = agent1.Require(agent2)
	registry := NewModuleRegistry([]Agent{})
	agent1.signUp(registry)
	result := registry.getAll()
	slices.SortFunc(result, func(a, b Agent) int {
		return cmp.Compare(a.moduleName(), b.moduleName())
	})
	assert.Equal([]Agent{agent1, agent2}, result)
}

func TestUpdate(t *testing.T) {
	assert := assertions.New(t)
	agent1Settings := ""
	agent2Settings := ""
	agent1Format := ""
	agent2Format := ""
	sequence := make([]string, 0)
	agent1 := NewAgent("1", func(r io.Reader, format string) error {
		bufReader := bufio.NewReader(r)
		agent1Settings, _ = bufReader.ReadString('\n')
		agent1Format = format
		sequence = append(sequence, "1")
		return nil
	})
	agent2 := NewAgent("2", func(r io.Reader, format string) error {
		bufReader := bufio.NewReader(r)
		agent2Settings, _ = bufReader.ReadString('\n')
		agent2Format = format
		sequence = append(sequence, "2")
		return nil
	})
	_ = agent1.Require(agent2)
	buffer := strings.NewReader("hello, world\n")
	reader := io.NewSectionReader(buffer, 0, buffer.Size())
	agent1.update(reader, "123")
	expectedSettings := "hello, world\n"
	expectedFormat := "123"
	assert.Equal(expectedSettings, agent1Settings)
	assert.Equal(expectedSettings, agent2Settings)
	assert.Equal(expectedFormat, agent1Format)
	assert.Equal(expectedFormat, agent2Format)
	assert.Equal([]string{"2", "1"}, sequence)
}

func TestUpdateChildReturnError(t *testing.T) {
	assert := assertions.New(t)
	now := time.Now()
	agent1Settings := ""
	agent2Settings := ""
	agent1Format := ""
	agent2Format := ""
	sequence := make([]string, 0)
	someError := errors.New("some error")
	agent1 := NewAgent("1", func(r io.Reader, format string) error {
		bufReader := bufio.NewReader(r)
		agent1Settings, _ = bufReader.ReadString('\n')
		agent1Format = format
		sequence = append(sequence, "1")
		return nil
	})
	agent2 := NewAgent("2", func(r io.Reader, format string) error {
		bufReader := bufio.NewReader(r)
		agent2Settings, _ = bufReader.ReadString('\n')
		agent2Format = format
		sequence = append(sequence, "2")
		return someError
	})
	_ = agent1.Require(agent2)
	buffer := strings.NewReader("hello, world\n")
	reader := io.NewSectionReader(buffer, 0, buffer.Size())
	err := agent1.update(reader, "123")
	expectedSettings := "hello, world\n"
	expectedFormat := "123"
	assert.Equal(someError, err)
	assert.Equal("", agent1Settings)
	assert.Equal(expectedSettings, agent2Settings)
	assert.Equal("", agent1Format)
	assert.Equal(expectedFormat, agent2Format)
	assert.Equal([]string{"2"}, sequence)
	assert.False(agent1.isConfigured(now))
	assert.False(agent2.isConfigured(now))
}

func TestUpdateRootReturnError(t *testing.T) {
	assert := assertions.New(t)
	now := time.Now()
	someError := errors.New("some error")
	agent1 := NewAgent("1", func(r io.Reader, format string) error {
		return someError
	})
	agent2 := NewAgent("2", func(r io.Reader, format string) error {
		return nil
	})
	_ = agent1.Require(agent2)
	buffer := strings.NewReader("hello, world\n")
	reader := io.NewSectionReader(buffer, 0, buffer.Size())
	err := agent1.update(reader, "123")
	assert.Equal(someError, err)
	assert.False(agent1.isConfigured(now))
	assert.True(agent2.isConfigured(now))
}
