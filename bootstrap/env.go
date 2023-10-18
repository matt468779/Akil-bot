package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	ServerAddress     string `mapstructure:"SERVER_ADDRESS"`
	AppEnv            string `mapstructure:"APP_ENV"`
	BaseURL           string `mapstructure:"BASE_URL"`
	OpenAIAPIKey      string `mapstructure:"OPENAI_API_KEY"`
	Port              string `mapstructure:"PORT"`
	BotToken          string `mapstructure:"BOT_TOKEN"`
	DBUri             string `mapstructure:"MONGODB_URI"`
	DBName            string `mapstructure:"DB_NAME"`
	ContextTimeout    int    `mapstructure:"CONTEXT_TIMEOUT"`
	ChannelID         int    `mapstructure:"CHANNEL_ID"`
	SystemMessagePath string `mapstructure:"SYSTEM_MESSAGE_PATH"`
	BackendURL        string `mapstructure:"BACKEND_URL"`
	FrontendURL       string `mapstructure:"FRONTEND_URL"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile("/etc/secrets/.env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
