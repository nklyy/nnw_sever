package config

import (
	"os"
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	type env struct {
		port            string
		prometheus      string
		environment     string
		mongoDbUrl      string
		mongoDbUser     string
		mongoDbPass     string
		mongoDbName     string
		jwtSecretKey    string
		shift           string
		passwordSalt    string
		emailFrom       string
		smtpHost        string
		smtpPort        string
		smtpUserApiKey  string
		smtpPasswordKey string
		twoFaIssuer     string
		devOrigin       string
		prodOrigin      string
	}

	type args struct {
		env env
	}

	setEnv := func(env env) {
		os.Setenv("PORT", env.port)
		os.Setenv("PROMETHEUS", env.prometheus)
		os.Setenv("ENVIRONMENT", env.environment)
		os.Setenv("MONGO_DB_NAME", env.mongoDbName)
		os.Setenv("MONGO_DB_USER", env.mongoDbUser)
		os.Setenv("MONGO_DB_PASS", env.mongoDbPass)
		os.Setenv("MONGO_DB_URL", env.mongoDbUrl)
		os.Setenv("JWT_SECRET_KEY", env.jwtSecretKey)
		os.Setenv("SHIFT", env.shift)
		os.Setenv("PASSWORD_SALT", env.passwordSalt)
		os.Setenv("EMAIL_FROM", env.emailFrom)
		os.Setenv("SMTP_HOST", env.smtpHost)
		os.Setenv("SMTP_PORT", env.smtpPort)
		os.Setenv("SMTP_USER_API_KEY", env.smtpUserApiKey)
		os.Setenv("SMTP_PASSWORD_KEY", env.smtpPasswordKey)
		os.Setenv("TWO_FA_ISSUER", env.twoFaIssuer)
		os.Setenv("DEV_ORIGIN", env.devOrigin)
		os.Setenv("PROD_ORIGIN", env.prodOrigin)
	}

	tests := []struct {
		name      string
		args      args
		want      *Config
		wantError bool
	}{
		{
			name: "Test config file!",
			args: args{
				env: env{
					port:            "4000",
					prometheus:      "9090",
					environment:     "development",
					mongoDbUrl:      "mongodb+srv://user:user@cluster0.database.mongodb.net/name?retryWrites=true&w=majority",
					mongoDbUser:     "admin",
					mongoDbPass:     "qwerty",
					mongoDbName:     "databaseName",
					jwtSecretKey:    "123qwerty",
					shift:           "123",
					passwordSalt:    "123",
					emailFrom:       "example@example.com",
					smtpHost:        "smtp.email.com",
					smtpPort:        "25",
					smtpUserApiKey:  "key",
					smtpPasswordKey: "password",
					twoFaIssuer:     "Example",
					devOrigin:       "http://localhost:3000",
					prodOrigin:      "https://example.com",
				},
			},
			want: &Config{
				PORT:        "4000",
				PROMETHEUS:  "9090",
				Environment: "development",
				EmailFrom:   "example@example.com",
				TwoFAIssuer: "Example",

				Secrets: Secrets{
					JwtSecretKey: "123qwerty",
					Shift:        123,
					PasswordSalt: 123,
				},

				MongoConfig: MongoConfig{
					MongoDbName: "databaseName",
					MongoDbUser: "admin",
					MongoDbPass: "qwerty",
					MongoDbUrl:  "mongodb+srv://user:user@cluster0.database.mongodb.net/name?retryWrites=true&w=majority",
				},

				SMTPConfig: SMTPConfig{
					SmtpHost:        "smtp.email.com",
					SmtpPort:        25,
					SmtpUserApiKey:  "key",
					SmtpPasswordKey: "password",
				},

				CorsOrigin: CorsOrigin{
					DevOrigin:  "http://localhost:3000",
					ProdOrigin: "https://example.com",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setEnv(test.args.env)

			got, err := Get()
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
