package credentials

import "github.com/pquerna/otp"

type Credentials struct {
	Password  string    `bson:"password"`
	SecretOTP SecretOTP `bson:"secret_otp"`
}

func (creds *Credentials) SetSecretOTP(key *otp.Key) {
	secretOtp := key.Secret()
	creds.SecretOTP = &secretOtp
}

type SecretOTP *string

var NilSecretOTP SecretOTP = nil
