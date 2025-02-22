package configurator

type Registry interface {
	get(key string) Agent
	set(key string, agent Agent)
	getAll() []Agent
}

func NewModuleRegistry(rootAgents []Agent) (Registry, error) {
	r := &moduleRegistry{
		agents: make(map[string]Agent),
	}
	for _, agent := range rootAgents {
		if err := agent.signUp(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

type moduleRegistry struct {
	agents map[string]Agent
}

func (r *moduleRegistry) get(key string) Agent {
	return r.agents[key]
}

func (r *moduleRegistry) set(key string, agent Agent) {
	r.agents[key] = agent
}

func (r *moduleRegistry) getAll() []Agent {
	list := make([]Agent, 0, len(r.agents))
	for _, value := range r.agents {
		list = append(list, value)
	}
	return list
}
