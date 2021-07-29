package config

type Configurations struct {
	PORT string `mapstructure:"PORT"`

	MongoDbName string `mapstructure:"MONGO_DB_NAME"`
	MongoDbUrl  string `mapstructure:"MONGO_DB_URL"`
}
