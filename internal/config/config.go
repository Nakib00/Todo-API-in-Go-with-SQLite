package config

type Config struct {
	ServerPort   string
	DatabasePath string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort:   "8080",
		DatabasePath: "./todos.db",
	}
}