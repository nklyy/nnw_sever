package config

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/spf13/viper"
	"os"
	"unicode"
)

type Configurations struct {
	PORT string `mapstructure:"PORT"`

	MongoDbName  string `mapstructure:"MONGO_DB_NAME"`
	MongoDbUser  string `mapstructure:"MONGO_DB_USER"`
	MongoDbPass  string `mapstructure:"MONGO_DB_PASS"`
	MongoDbUrl   string `mapstructure:"MONGO_DB_URL"`
	JwtSecretKey string `mapstructure:"JWT_SECRET_KEY"`
	Shift        string `mapstructure:"SHIFT"`
	PasswordSalt string `mapstructure:"PASSWORD_SALT"`
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
}

func ValidatorConfig(v *validator.Validate) ut.Translator {
	translator := en.New()
	uni := ut.New(translator, translator)

	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, found := uni.GetTranslator("en")
	if !found {
		fmt.Printf("ERROR: %s \n", "translator not found")
		return nil
	}

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return nil
	}

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("passwd", trans, func(ut ut.Translator) error {
		return ut.Add("passwd", "{0} is not strong enough", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("passwd", fe.Field())
		return t
	})

	_ = v.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {
		var (
			hasMinLen = false
			hasUpper  = false
			hasLower  = false
			hasNumber = false
		)
		if len(fl.Field().String()) > 7 {
			hasMinLen = true
		}
		for _, char := range fl.Field().String() {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsNumber(char):
				hasNumber = true
			}
		}
		return hasMinLen && hasUpper && hasLower && hasNumber
	})

	return trans
}
