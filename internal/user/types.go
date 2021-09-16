package user

type (
	Status    string
	SecretOTP *string
)

const (
	Active   Status = "active"
	Disabled Status = "disabled"
)

var NilSecretOTP SecretOTP = nil
