package config

import "github.com/spf13/viper"

type Configurations struct {
	PORT string `mapstructure:"PORT"`

	MongoDbName  string `mapstructure:"MONGO_DB_NAME"`
	MongoDbUser  string `mapstructure:"MONGO_DB_USER"`
	MongoDbPass  string `mapstructure:"MONGO_DB_PASS"`
	MongoDbUrl   string `mapstructure:"MONGO_DB_URL"`
	JwtSecretKey string `mapstructure:"JWT_SECRET_KEY"`
	Shift        string `mapstructure:"SHIFT"`
}

func InitConfig(path string) (*Configurations, error) {
	viper.AddConfigPath(path)

	viper.SetConfigName("app")

	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var configuration Configurations
	err = viper.Unmarshal(&configuration)
	if err != nil {
		//fmt.Printf("Unable to decode into struct, %v", err)
		return nil, err
	}

	return &configuration, nil
}
