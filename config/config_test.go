package config

import (
	"os"
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	type env struct {
		port            string
		mongoDbUrl      string
		mongoDbUser     string
		mongoDbPass     string
		mongoDbName     string
		jwtSecretKey    string
		shift           string
		passwordSalt    string
		smtpHost        string
		smtpPort        string
		smtpUserApiKey  string
		smtpPasswordKey string
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
		os.Setenv("JWT_SECRET_KEY", env.jwtSecretKey)
		os.Setenv("SHIFT", env.shift)
		os.Setenv("PASSWORD_SALT", env.passwordSalt)
		os.Setenv("SMTP_HOST", env.smtpHost)
		os.Setenv("SMTP_PORT", env.smtpPort)
		os.Setenv("SMTP_USER_API_KEY", env.smtpUserApiKey)
		os.Setenv("SMTP_PASSWORD_KEY", env.smtpPasswordKey)
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
					port:            ":4000",
					mongoDbName:     "databaseName",
					mongoDbUser:     "admin",
					mongoDbPass:     "qwerty",
					mongoDbUrl:      "mongodb+srv://user:user@cluster0.database.mongodb.net/name?retryWrites=true&w=majority",
					jwtSecretKey:    "123qwerty",
					shift:           "123",
					passwordSalt:    "123",
					smtpHost:        "smtp.email.com",
					smtpPort:        "25",
					smtpUserApiKey:  "key",
					smtpPasswordKey: "password",
				},
				path: "..",
			},
			want: &Configurations{
				PORT:            ":4000",
				MongoDbName:     "databaseName",
				MongoDbUser:     "admin",
				MongoDbPass:     "qwerty",
				MongoDbUrl:      "mongodb+srv://user:user@cluster0.database.mongodb.net/name?retryWrites=true&w=majority",
				JwtSecretKey:    "123qwerty",
				Shift:           "123",
				PasswordSalt:    "123",
				SmtpHost:        "smtp.email.com",
				SmtpPort:        "25",
				SmtpUserApiKey:  "key",
				SmtpPasswordKey: "password",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setEnv(test.args.env)

			got, err := InitConfig(test.args.path, "")
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
