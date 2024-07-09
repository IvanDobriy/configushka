package configurator

type Registry interface {
	Get(key string) Agent
	Set(key string, agent Agent)
	GetAll() []Agent
}

var registry = map[string]Agent{}

func NewRegistry() Registry {
	r := &RegistryImpl{}
	return r
}

type RegistryImpl struct {
}

func (r *RegistryImpl) Get(key string) Agent {
	return registry[key]
}

func (r *RegistryImpl) Set(key string, agent Agent) {
	registry[key] = agent
}

func (r *RegistryImpl) GetAll() []Agent {
	list := make([]Agent, 0, len(registry))
	for _, value := range registry {
		list = append(list, value)
	}
	return list
}
