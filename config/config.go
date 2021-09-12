package config

import (
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PORT        string `default:"4000" envconfig:"PORT"`
	Environment string `default:"development" envconfig:"APP_ENV"`
	EmailFrom   string `envconfig:"EMAIL_FROM"`

	Secrets
	MongoConfig
	SMTPConfig
}

type Secrets struct {
	JwtSecretKey string `envconfig:"JWT_SECRET_KEY"`
	Shift        int    `envconfig:"SHIFT"`
	PasswordSalt int    `envconfig:"PASSWORD_SALT"`
}

type MongoConfig struct {
	MongoDbName string `envconfig:"MONGO_DB_NAME"`
	MongoDbUser string `envconfig:"MONGO_DB_USER"`
	MongoDbPass string `envconfig:"MONGO_DB_PASS"`
	MongoDbUrl  string `envconfig:"MONGO_DB_URL"`
}

type SMTPConfig struct {
	SmtpHost        string `envconfig:"SMTP_HOST"`
	SmtpPort        string `envconfig:"SMTP_PORT"`
	SmtpUserApiKey  string `envconfig:"SMTP_USER_API_KEY"`
	SmtpPasswordKey string `envconfig:"SMTP_PASSWORD_KEY"`
}

var (
	once   sync.Once
	config *Config
)

func Get() (*Config, error) {
	var err error
	once.Do(func() {
		var cfg Config
		_ = godotenv.Load("../.env")

		if err = envconfig.Process("", &cfg); err != nil {
			return
		}

		config = &cfg
	})

	return config, err
}
