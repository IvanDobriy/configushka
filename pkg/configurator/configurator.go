package configurator

import "strings"

type Configurator interface {
	Configure() error
}

func NewConfigurator(rootAgents []Agent) Configurator {
	configurator := &ConfiguratorImpl{
		rootAgents: rootAgents,
	}
	return configurator
}

type ConfiguratorImpl struct {
	rootAgents []Agent
}

func (c *ConfiguratorImpl) Configure() error {
	r := NewRegistry()
	registeredAgents := r.GetAll()
	conf := strings.NewReader("hello, world")
	for _, agent := range registeredAgents {
		conf.Seek(0, 0)
		agent.Update(conf)
	}
	return nil
}

func (c *ConfiguratorImpl) update(agent *Agent) error {
	return nil
}
