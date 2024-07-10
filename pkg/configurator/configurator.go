package configurator

import (
	"errors"
	"os"
)

type Configurator interface {
	Configure() error
}

func NewConfigurator(registry Registry, configPaths []string) Configurator {
	configurator := &configuratorImpl{
		registry: registry,
		paths:    configPaths,
	}
	return configurator
}

type configuratorImpl struct {
	registry Registry
	paths    []string
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
		conf.Seek(0, 0)
		agent.update(conf)
	}
	return nil
}
