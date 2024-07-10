package configurator

import (
	"errors"
	"os"
)

type Configurator interface {
	Configure() error
}

func NewConfigurator(registry Registry, configPaths []string, format string) Configurator {
	configurator := &configuratorImpl{
		registry: registry,
		paths:    configPaths,
		format:   format,
	}
	return configurator
}

type configuratorImpl struct {
	registry Registry
	paths    []string
	format   string
}

func (c *configuratorImpl) Configure() error {
	registeredAgents := c.registry.getAll()
	if len(c.paths) == 0 {
		return errors.New("no config files found")
	}
	conf, error := os.OpenFile(c.paths[0], os.O_RDONLY, 0640)
	if error != nil {
		return error
	}
	for _, agent := range registeredAgents {
		agent.update(conf, c.format)
	}
	return nil
}
