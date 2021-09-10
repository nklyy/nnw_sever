package user

type CreateUserDTO struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	SecretOTP string `json:"secret_otp"`
}
