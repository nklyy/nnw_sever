package credentials

func MapToEntity(dto *DTO) *Credentials {
	return &Credentials{
		Password:  dto.Password,
		SecretOTP: dto.SecretOTP,
	}
}

func MapToDTO(credentials *Credentials) *DTO {
	return &DTO{
		Password:  credentials.Password,
		SecretOTP: credentials.SecretOTP,
	}
}
