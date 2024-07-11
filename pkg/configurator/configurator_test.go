package configurator

import (
	assertions "github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConfigure(t *testing.T) {
	assert := assertions.New(t)
	path, err := filepath.Abs("../../test/configurator/test.config.yaml")
	assert.Nil(err)
	now := time.Now()
	raw, err := os.ReadFile(path)
	assert.Nil(err)
	expected := string(raw)
	agent1Config := ""
	agent1 := NewAgent("1", func(r io.Reader, format string) error {
		buffer, ioerr := io.ReadAll(r)
		agent1Config = string(buffer)
		return ioerr
	})
	registry, err := NewModuleRegistry([]Agent{agent1})
	assert.Nil(err)
	configurator := NewLocalConfigurator(registry, []string{path}, "yaml")

	err = configurator.Configure()
	assert.Nil(err)

	assert.Equal(expected, agent1Config)
	assert.True(agent1.isConfigured(now))
}

func TestConfigureSecondPathExists(t *testing.T) {
	assert := assertions.New(t)
	path, err := filepath.Abs("../../test/configurator/test.config.yaml")
	assert.Nil(err)
	wrongPath, err := filepath.Abs("../../test/configurator/wrongPath.config.yaml")
	assert.Nil(err)
	now := time.Now()
	raw, err := os.ReadFile(path)
	assert.Nil(err)
	expected := string(raw)
	agent1Config := ""
	agent1 := NewAgent("1", func(r io.Reader, format string) error {
		buffer, ioerr := io.ReadAll(r)
		agent1Config = string(buffer)
		return ioerr
	})
	registry, err := NewModuleRegistry([]Agent{agent1})
	assert.Nil(err)
	configurator := NewLocalConfigurator(registry, []string{wrongPath, path}, "yaml")

	err = configurator.Configure()
	assert.Nil(err)

	assert.Equal(expected, agent1Config)
	assert.True(agent1.isConfigured(now))
}

func TestConfigureFileNotFound(t *testing.T) {
	assert := assertions.New(t)
	wrongPath, err := filepath.Abs("../../test/configurator/wrongPath.config.yaml")
	assert.Nil(err)
	now := time.Now()
	agent1Config := ""
	agent1 := NewAgent("1", func(r io.Reader, format string) error {
		buffer, ioerr := io.ReadAll(r)
		agent1Config = string(buffer)
		return ioerr
	})
	registry, err := NewModuleRegistry([]Agent{agent1})
	assert.Nil(err)
	configurator := NewLocalConfigurator(registry, []string{wrongPath}, "yaml")
	err = configurator.Configure()
	assert.NotNil(err)
	assert.Equal("", agent1Config)
	assert.False(agent1.isConfigured(now))
}

func TestConfigureFormatAlwaysLowerCase(t *testing.T) {
	assert := assertions.New(t)
	path, err := filepath.Abs("../../test/configurator/test.config.yaml")
	assert.Nil(err)
	now := time.Now()
	agent1Format := ""
	agent1 := NewAgent("1", func(r io.Reader, format string) error {
		agent1Format = format
		return nil
	})
	registry, err := NewModuleRegistry([]Agent{agent1})
	assert.Nil(err)
	configurator := NewLocalConfigurator(registry, []string{path}, "yAmL")

	err = configurator.Configure()
	assert.Nil(err)
	assert.Equal("yaml", agent1Format)
	assert.True(agent1.isConfigured(now))
}

func TestConfigureTwoAgents(t *testing.T) {
	assert := assertions.New(t)
	path, err := filepath.Abs("../../test/configurator/test.config.yaml")
	assert.Nil(err)
	now := time.Now()

	agent1 := NewAgent("1", func(r io.Reader, format string) error {
		return nil
	})
	agent2 := NewAgent("2", func(r io.Reader, format string) error {
		return nil
	})

	registry, err := NewModuleRegistry([]Agent{agent1, agent2})
	assert.Nil(err)
	configurator := NewLocalConfigurator(registry, []string{path}, "yaml")

	err = configurator.Configure()
	assert.Nil(err)

	assert.True(agent1.isConfigured(now))
	assert.True(agent2.isConfigured(now))
}

func TestConfigure2LevelHierarchy(t *testing.T) {
	assert := assertions.New(t)
	path, err := filepath.Abs("../../test/configurator/test.config.yaml")
	assert.Nil(err)
	now := time.Now()
	sequence := make([]string, 0)
	agent1 := NewAgent("1", func(r io.Reader, format string) error {
		sequence = append(sequence, "1")
		return nil
	})
	agent2 := NewAgent("2", func(r io.Reader, format string) error {
		sequence = append(sequence, "2")
		return nil
	})
	agent1.Require(agent2)

	registry, err := NewModuleRegistry([]Agent{agent1, agent2})
	assert.Nil(err)
	configurator := NewLocalConfigurator(registry, []string{path}, "yaml")

	err = configurator.Configure()
	assert.Nil(err)
	assert.Equal([]string{"2", "1"}, sequence)
	assert.True(agent1.isConfigured(now))
	assert.True(agent2.isConfigured(now))
}
