package configurator

import "strings"

type Configurator interface {
	Configure() error
}

func NewConfigurator(registry Registry) Configurator {
	configurator := &configuratorImpl{
		registry: registry,
	}
	return configurator
}

type configuratorImpl struct {
	registry Registry
}

func (c *configuratorImpl) Configure() error {
	registeredAgents := c.registry.getAll()
	conf := strings.NewReader("hello, world")
	for _, agent := range registeredAgents {
		conf.Seek(0, 0)
		agent.update(conf)
	}
	return nil
}
