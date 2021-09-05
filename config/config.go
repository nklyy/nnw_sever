package config

import (
	"github.com/spf13/viper"
	"os"
)

type Configurations struct {
	PORT string `mapstructure:"PORT"`

	MongoDbName string `mapstructure:"MONGO_DB_NAME"`
	MongoDbUser string `mapstructure:"MONGO_DB_USER"`
	MongoDbPass string `mapstructure:"MONGO_DB_PASS"`
	MongoDbUrl  string `mapstructure:"MONGO_DB_URL"`

	JwtSecretKey string `mapstructure:"JWT_SECRET_KEY"`

	Shift        string `mapstructure:"SHIFT"`
	PasswordSalt string `mapstructure:"PASSWORD_SALT"`

	EmailFrom       string `mapstructure:"EMAIL_FROM"`
	SmtpHost        string `mapstructure:"SMTP_HOST"`
	SmtpPort        string `mapstructure:"SMTP_PORT"`
	SmtpUserApiKey  string `mapstructure:"SMTP_USER_API_KEY"`
	SmtpPasswordKey string `mapstructure:"SMTP_PASSWORD_KEY"`
}

func InitConfig(path string, env string) (*Configurations, error) {
	var configuration Configurations

	if env == "PRODUCTION" {
		setFromEnv(&configuration)
		return &configuration, nil
	}

	viper.AddConfigPath(path)

	viper.SetConfigName(".env")

	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&configuration)
	if err != nil {
		//fmt.Printf("Unable to decode into struct, %v", err)
		return nil, err
	}

	return &configuration, nil
}

func setFromEnv(cfg *Configurations) {
	cfg.MongoDbName = os.Getenv("MONGO_DB_NAME")
	cfg.MongoDbUser = os.Getenv("MONGO_DB_USER")
	cfg.MongoDbPass = os.Getenv("MONGO_DB_PASS")
	cfg.MongoDbUrl = os.Getenv("MONGO_DB_URL")

	cfg.JwtSecretKey = os.Getenv("JWT_SECRET_KEY")

	cfg.Shift = os.Getenv("SHIFT")
	cfg.PasswordSalt = os.Getenv("PASSWORD_SALT")

	cfg.PORT = os.Getenv("PORT")

	cfg.EmailFrom = os.Getenv("EMAIL_FROM")
	cfg.SmtpHost = os.Getenv("SMTP_HOST")
	cfg.SmtpPort = os.Getenv("SMTP_PORT")
	cfg.SmtpUserApiKey = os.Getenv("SMTP_USER_API_KEY")
	cfg.SmtpPasswordKey = os.Getenv("SMTP_PASSWORD_KEY")
}
