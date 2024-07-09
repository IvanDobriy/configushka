package configurator

type Configurator interface {
	Configure(agents []Agent) error
}

func NewConfigurator() Configurator {
	configurator := &ConfiguratorImpl{}
	return configurator
}

type ConfiguratorImpl struct {
}

func (c *ConfiguratorImpl) Configure(agents []Agent) error {
	return nil
}

func (c *ConfiguratorImpl) update(agent *Agent) error {
	return nil
}
