package configurator

import (
	"errors"
	"os"
	"strings"
)

type Configurator interface {
	Configure() error
}

func NewLocalConfigurator(registry Registry, configPaths []string, format string) Configurator {
	configurator := &configuratorImpl{
		registry: registry,
		paths:    configPaths,
		format:   strings.ToLower(format),
	}
	return configurator
}

type configuratorImpl struct {
	registry Registry
	paths    []string
	format   string
}

func (c *configuratorImpl) Configure() (err error) {
	registeredAgents := c.registry.getAll()
	if len(c.paths) == 0 {
		return errors.New("configuration file paths is empty")
	}
	var conf *os.File = nil
	for _, path := range c.paths {
		file, error := os.OpenFile(path, os.O_RDONLY, 0640)
		if error == nil {
			conf = file
			break
		}
	}
	if conf == nil {
		return errors.New("no config files found")
	}

	defer func() {
		err = conf.Close()
	}()
	for _, agent := range registeredAgents {
		if err = agent.update(conf, c.format); err != nil {
			return
		}
	}
	return
}
