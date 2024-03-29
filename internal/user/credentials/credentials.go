package credentials

import (
	"github.com/pquerna/otp"
)

type Credentials struct {
	Password  string    `bson:"password"`
	SecretOTP SecretOTP `bson:"secret_otp"`
}

func (credentials *Credentials) SetSecretOTP(key *otp.Key) {
	secretOtp := key.Secret()
	credentials.SecretOTP = &secretOtp
}

func (credentials *Credentials) SetNewPassword(password string) {
	credentials.Password = password
}

type SecretOTP *string

var NilSecretOTP SecretOTP = nil
