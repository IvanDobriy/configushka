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
	registry := NewModuleRegistry([]Agent{agent1})
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
	registry := NewModuleRegistry([]Agent{agent1})
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
	registry := NewModuleRegistry([]Agent{agent1})
	configurator := NewLocalConfigurator(registry, []string{wrongPath}, "yaml")
	err = configurator.Configure()
	assert.NotNil(err)
	assert.Equal("", agent1Config)
	assert.False(agent1.isConfigured(now))
}
