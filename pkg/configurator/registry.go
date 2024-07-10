package configurator

type registry interface {
	get(key string) Agent
	set(key string, agent Agent)
	getAll() []Agent
}

var __registry = map[string]Agent{}

func getModuleRegistry() registry {
	r := &moduleRegistry{}
	return r
}

type moduleRegistry struct {
}

func (r *moduleRegistry) get(key string) Agent {
	return __registry[key]
}

func (r *moduleRegistry) set(key string, agent Agent) {
	__registry[key] = agent
}

func (r *moduleRegistry) getAll() []Agent {
	list := make([]Agent, 0, len(__registry))
	for _, value := range __registry {
		list = append(list, value)
	}
	return list
}
