package configurator

import "strings"

type Configurator interface {
	Configure() error
}

func NewConfigurator(rootAgents []Agent) Configurator {
	configurator := &ConfiguratorImpl{}
	for _, agent := range rootAgents {
		agent.register()
	}
	return configurator
}

type ConfiguratorImpl struct {
}

func (c *ConfiguratorImpl) Configure() error {
	r := getModuleRegistry()
	registeredAgents := r.getAll()
	conf := strings.NewReader("hello, world")
	for _, agent := range registeredAgents {
		conf.Seek(0, 0)
		agent.update(conf)
	}
	return nil
}
