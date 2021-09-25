package config

import (
	"encoding/json"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PORT        string `required:"true" default:"4000" envconfig:"PORT"`
	Environment string `required:"true" default:"development" envconfig:"APP_ENV"`
	EmailFrom   string `required:"true" envconfig:"EMAIL_FROM"`
	TwoFAIssuer string `required:"true" envconfig:"TWO_FA_ISSUER" default:"NNW"`

	Secrets
	MongoConfig
	SMTPConfig
	CorsOrigin
}

func (cfg Config) String() string {
	buf, _ := json.MarshalIndent(&cfg, "", "	")
	return string(buf)
}

type Secrets struct {
	JwtSecretKey string `required:"true" envconfig:"JWT_SECRET_KEY"`
	Shift        int    `required:"true" envconfig:"SHIFT"`
	PasswordSalt int    `required:"true" envconfig:"PASSWORD_SALT"`
}

type MongoConfig struct {
	MongoDbName string `required:"true" envconfig:"MONGO_DB_NAME"`
	MongoDbUser string `required:"true" envconfig:"MONGO_DB_USER"`
	MongoDbPass string `required:"true" envconfig:"MONGO_DB_PASS"`
	MongoDbUrl  string `required:"true" envconfig:"MONGO_DB_URL"`
}

type SMTPConfig struct {
	SmtpHost        string `required:"true" envconfig:"SMTP_HOST"`
	SmtpPort        int    `required:"true" envconfig:"SMTP_PORT"`
	SmtpUserApiKey  string `required:"true" envconfig:"SMTP_USER_API_KEY"`
	SmtpPasswordKey string `required:"true" envconfig:"SMTP_PASSWORD_KEY"`
}

type CorsOrigin struct {
	DevOrigin  string `required:"true" envconfig:"DEV_ORIGIN"`
	ProdOrigin string `required:"true" envconfig:"PROD_ORIGIN"`
}

var (
	once   sync.Once
	config *Config
)

func Get() (*Config, error) {
	var err error
	once.Do(func() {
		var cfg Config
		// If you run it locally and through terminal please set up this in Load function (../.env)
		_ = godotenv.Load(".env")

		if err = envconfig.Process("", &cfg); err != nil {
			return
		}

		config = &cfg
	})

	return config, err
}
