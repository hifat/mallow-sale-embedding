package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type App struct {
	Name string
	Host string
	Port int
	Mode string
}

type Auth struct {
	ApiKey string
}

type QDB struct {
	Host   string
	Port   int
	ApiKey string
}

type Agent struct {
	OllamaHost string
}

type Config struct {
	App   App
	QDB   QDB
	Auth  Auth
	Agent Agent
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = ".env"
	}

	viper.SetConfigFile(path)
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found, using environment variables")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfg := Config{
		App: App{
			Port: viper.GetInt("APP_PORT"),
		},
		Auth: Auth{
			ApiKey: viper.GetString("API_KEY"),
		},
		QDB: QDB{
			Host:   viper.GetString("QD_HOST"),
			Port:   viper.GetInt("QD_PORT"),
			ApiKey: viper.GetString("QD_API_KEY"),
		},
		Agent: Agent{
			OllamaHost: viper.GetString("OLLAMA_HOST"),
		},
	}

	return &cfg, nil
}
