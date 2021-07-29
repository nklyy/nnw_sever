package config

type Configurations struct {
	PORT string `mapstructure:"PORT"`

	MongoDbUrl string `mapstructure:"MONGO_DB_URL"`
}
