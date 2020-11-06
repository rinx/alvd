package config

type Config struct {
	AgentEnabled bool
}

func New() (*Config, error) {
	return &Config{}, nil
}
