package config

import (
	"os"
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	type env struct {
		port        string
		mongoDbUrl  string
		mongoDbUser string
		mongoDbPass string
		mongoDbName string
	}

	type args struct {
		path string
		env  env
	}

	setEnv := func(env env) {
		os.Setenv("PORT", env.port)
		os.Setenv("MONGO_DB_NAME", env.mongoDbName)
		os.Setenv("MONGO_DB_USER", env.mongoDbUser)
		os.Setenv("MONGO_DB_PASS", env.mongoDbPass)
		os.Setenv("MONGO_DB_URL", env.mongoDbUrl)
	}

	tests := []struct {
		name      string
		args      args
		want      *Configurations
		wantError bool
	}{
		{
			name: "Test config file!",
			args: args{
				env: env{
					port:        ":4000",
					mongoDbName: "databaseName",
					mongoDbUser: "admin",
					mongoDbPass: "qwerty",
					mongoDbUrl:  "mongodb+srv://user:user@cluster0.database.mongodb.net/name?retryWrites=true&w=majority",
				},
				path: ".",
			},
			want: &Configurations{
				PORT:        ":4000",
				MongoDbName: "databaseName",
				MongoDbUser: "admin",
				MongoDbPass: "qwerty",
				MongoDbUrl:  "mongodb+srv://user:user@cluster0.database.mongodb.net/name?retryWrites=true&w=majority",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setEnv(test.args.env)

			got, err := InitConfig(test.args.path)
			if (err != nil) != test.wantError {
				t.Errorf("Init() error = %v, wantErr %v", err, test.wantError)

				return
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Init() got = %v, want %v", got, test.want)
			}
		})
	}
}
