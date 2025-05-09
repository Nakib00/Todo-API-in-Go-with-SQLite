package config

type Config struct {
	ServerPort string
	DBConfig   struct {
		Username string
		Password string
		Host     string
		Port     string
		DBName   string
	}
}

func LoadConfig() *Config {
	cfg := &Config{
		ServerPort: "8080",
	}
	cfg.DBConfig.Username = "root"
	cfg.DBConfig.Password = ""
	cfg.DBConfig.Host = "localhost"
	cfg.DBConfig.Port = "3306"
	cfg.DBConfig.DBName = "todo-go"
	return cfg
}
