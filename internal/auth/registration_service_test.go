package auth

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	mock_jwt "nnw_s/internal/auth/jwt/mocks"
	mock_twofa "nnw_s/internal/auth/twofa/mocks"
	mock_verification "nnw_s/internal/auth/verification/mocks"
	"nnw_s/internal/user"
	"nnw_s/internal/user/credentials"
	mock_credentials "nnw_s/internal/user/credentials/mocks"
	mock_user "nnw_s/internal/user/mocks"
	"nnw_s/pkg/errors"
	"nnw_s/pkg/notificator"
	mock_notificator "nnw_s/pkg/notificator/mocks"
	"nnw_s/pkg/wallet"
	"testing"
)

func TestNewRegistrationService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	emailSender := "example@example.com"

	tests := []struct {
		name        string
		log         *logrus.Logger
		deps        *ServiceDeps
		emailSender string
		expect      func(*testing.T, RegistrationService, error)
	}{
		{
			name: "should return registration service",
			log:  logrus.New(),
			deps: &ServiceDeps{
				UserService:         mock_user.NewMockService(controller),
				NotificatorService:  mock_notificator.NewMockService(controller),
				VerificationService: mock_verification.NewMockService(controller),
				TwoFAService:        mock_twofa.NewMockService(controller),
				JWTService:          mock_jwt.NewMockService(controller),
				CredentialsService:  mock_credentials.NewMockService(controller),
			},
			emailSender: emailSender,
			expect: func(t *testing.T, service RegistrationService, err error) {
				assert.NotNil(t, service)
				assert.Nil(t, err)
			},
		},
		{
			name:        "should return invalid service dependencies",
			log:         logrus.New(),
			deps:        nil,
			emailSender: emailSender,
			expect: func(t *testing.T, service RegistrationService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid service dependencies")
			},
		},
		{
			name: "should return invalid user service",
			log:  logrus.New(),
			deps: &ServiceDeps{
				UserService:         nil,
				NotificatorService:  mock_notificator.NewMockService(controller),
				VerificationService: mock_verification.NewMockService(controller),
				TwoFAService:        mock_twofa.NewMockService(controller),
				JWTService:          mock_jwt.NewMockService(controller),
				CredentialsService:  mock_credentials.NewMockService(controller),
			},
			emailSender: emailSender,
			expect: func(t *testing.T, service RegistrationService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid user service")
			},
		},
		{
			name: "should return invalid notification service",
			log:  logrus.New(),
			deps: &ServiceDeps{
				UserService:         mock_user.NewMockService(controller),
				NotificatorService:  nil,
				VerificationService: mock_verification.NewMockService(controller),
				TwoFAService:        mock_twofa.NewMockService(controller),
				JWTService:          mock_jwt.NewMockService(controller),
				CredentialsService:  mock_credentials.NewMockService(controller),
			},
			emailSender: emailSender,
			expect: func(t *testing.T, service RegistrationService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid notification service")
			},
		},
		{
			name: "should return invalid verification service",
			log:  logrus.New(),
			deps: &ServiceDeps{
				UserService:         mock_user.NewMockService(controller),
				NotificatorService:  mock_notificator.NewMockService(controller),
				VerificationService: nil,
				TwoFAService:        mock_twofa.NewMockService(controller),
				JWTService:          mock_jwt.NewMockService(controller),
				CredentialsService:  mock_credentials.NewMockService(controller),
			},
			emailSender: emailSender,
			expect: func(t *testing.T, service RegistrationService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid verification service")
			},
		},
		{
			name: "should return invalid TwoFA service",
			log:  logrus.New(),
			deps: &ServiceDeps{
				UserService:         mock_user.NewMockService(controller),
				NotificatorService:  mock_notificator.NewMockService(controller),
				VerificationService: mock_verification.NewMockService(controller),
				TwoFAService:        nil,
				JWTService:          mock_jwt.NewMockService(controller),
				CredentialsService:  mock_credentials.NewMockService(controller),
			},
			emailSender: emailSender,
			expect: func(t *testing.T, service RegistrationService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid TwoFA service")
			},
		},
		{
			name: "should return invalid JWT service",
			log:  logrus.New(),
			deps: &ServiceDeps{
				UserService:         mock_user.NewMockService(controller),
				NotificatorService:  mock_notificator.NewMockService(controller),
				VerificationService: mock_verification.NewMockService(controller),
				TwoFAService:        mock_twofa.NewMockService(controller),
				JWTService:          nil,
				CredentialsService:  mock_credentials.NewMockService(controller),
			},
			emailSender: emailSender,
			expect: func(t *testing.T, service RegistrationService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid JWT service")
			},
		},
		{
			name: "should return invalid credentials service",
			log:  logrus.New(),
			deps: &ServiceDeps{
				UserService:         mock_user.NewMockService(controller),
				NotificatorService:  mock_notificator.NewMockService(controller),
				VerificationService: mock_verification.NewMockService(controller),
				TwoFAService:        mock_twofa.NewMockService(controller),
				JWTService:          mock_jwt.NewMockService(controller),
				CredentialsService:  nil,
			},
			emailSender: emailSender,
			expect: func(t *testing.T, service RegistrationService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid credentials service")
			},
		},
		{
			name: "should return invalid logger",
			log:  nil,
			deps: &ServiceDeps{
				UserService:         mock_user.NewMockService(controller),
				NotificatorService:  mock_notificator.NewMockService(controller),
				VerificationService: mock_verification.NewMockService(controller),
				TwoFAService:        mock_twofa.NewMockService(controller),
				JWTService:          mock_jwt.NewMockService(controller),
				CredentialsService:  mock_credentials.NewMockService(controller),
			},
			emailSender: emailSender,
			expect: func(t *testing.T, service RegistrationService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid logger")
			},
		},
		{
			name: "should return invalid sender's email",
			log:  logrus.New(),
			deps: &ServiceDeps{
				UserService:         mock_user.NewMockService(controller),
				NotificatorService:  mock_notificator.NewMockService(controller),
				VerificationService: mock_verification.NewMockService(controller),
				TwoFAService:        mock_twofa.NewMockService(controller),
				JWTService:          mock_jwt.NewMockService(controller),
				CredentialsService:  mock_credentials.NewMockService(controller),
			},
			emailSender: "",
			expect: func(t *testing.T, service RegistrationService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid sender's email")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := NewRegistrationService(tc.log, tc.emailSender, tc.deps)
			tc.expect(t, svc, err)
		})
	}
}

func TestRegistrationSvc_RegisterUser(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	log := logrus.New()

	mockUserSvc := mock_user.NewMockService(controller)
	mockVerificationSvc := mock_verification.NewMockService(controller)
	mockNotificationSvc := mock_notificator.NewMockService(controller)

	deps := &ServiceDeps{
		UserService:         mockUserSvc,
		NotificatorService:  mockNotificationSvc,
		VerificationService: mockVerificationSvc,
		TwoFAService:        mock_twofa.NewMockService(controller),
		JWTService:          mock_jwt.NewMockService(controller),
		CredentialsService:  mock_credentials.NewMockService(controller),
	}

	// Test Data
	emailSender := "example@example.com"
	userEmail := "user@example.com"
	userPassword := "==WvZitmZDgzSHgAWvKs"
	code := "ASDDSA"

	service, _ := NewRegistrationService(log, emailSender, deps)

	var registerUserDTO RegisterUserDTO
	registerUserDTO.Email = userEmail
	registerUserDTO.Password = userPassword

	// Test DTO
	var setupNewPasswordDTO SetupNewPasswordDTO
	setupNewPasswordDTO.Email = userEmail
	setupNewPasswordDTO.Password = userPassword

	// Test Cred
	secretKey := "secret"
	var testCred credentials.Credentials
	testCred.Password = "==WvZitmZDgzSHgAWvKs"
	testCred.SecretOTP = &secretKey

	//testCredDTO := credentials.MapToDTO(&testCred)

	// Test wallet
	var testWallet []*wallet.Wallet
	testWallet = append(testWallet, &wallet.Wallet{
		Name:     "BTC",
		WalletId: "8ebdfa95-484d-11ec-ba92-38d547b6cf94",
		Address:  "mrgZBqLCicXRGfSjqiSiV39mXgsV3euVZt",
	})

	// Test user
	testUser, _ := user.NewUser(userEmail, &testWallet, &testCred)

	notActiveUser := user.MapToDTO(testUser)

	testUser.SetToActive()
	testUser.SetToVerified()

	testUserDTO := user.MapToDTO(testUser)

	wrongUserDTO := user.MapToDTO(testUser)
	wrongUserDTO.ID = "example"

	// Test Email data
	testEmailData := notificator.Email{
		Subject:   emailVerificationSubject,
		Recipient: userEmail,
		Sender:    emailSender,
		Template:  emailVerificationTemplateName,
		Data: map[string]interface{}{
			"topic":   emailVerificationTopic,
			"message": emailVerificationMessage,
			"code":    code,
		},
	}

	tests := []struct {
		name   string
		ctx    context.Context
		dto    *RegisterUserDTO
		setup  func(context.Context, *RegisterUserDTO)
		expect func(*testing.T, error)
	}{
		{
			name: "should return failed to register doesn't exist user",
			ctx:  context.Background(),
			dto:  &registerUserDTO,
			setup: func(ctx context.Context, dto *RegisterUserDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(nil, nil)
				mockUserSvc.EXPECT().CreateUser(ctx, &user.CreateUserDTO{
					Email:    dto.Email,
					Password: dto.Password,
				}).Return("", errors.NewInternal("Failed to create user"))
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("Failed to create user"), err)
			},
		},
		{
			name: "should return failed to delete user",
			ctx:  context.Background(),
			dto:  &registerUserDTO,
			setup: func(ctx context.Context, dto *RegisterUserDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(notActiveUser, nil)
				mockUserSvc.EXPECT().DeleteUserByEmail(ctx, dto.Email).Return(errors.NewInternal("Failed to delete user"))
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("Failed to delete user"), err)
			},
		},
		{
			name: "should return failed to register exist user",
			ctx:  context.Background(),
			dto:  &registerUserDTO,
			setup: func(ctx context.Context, dto *RegisterUserDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(notActiveUser, nil)
				mockUserSvc.EXPECT().DeleteUserByEmail(ctx, dto.Email).Return(nil)
				mockUserSvc.EXPECT().CreateUser(ctx, &user.CreateUserDTO{
					Email:    dto.Email,
					Password: dto.Password,
				}).Return("", errors.NewInternal("Failed to create user"))
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("Failed to create user"), err)
			},
		},
		{
			name: "should return user already exist",
			ctx:  context.Background(),
			dto:  &registerUserDTO,
			setup: func(ctx context.Context, dto *RegisterUserDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(testUserDTO, nil)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, user.ErrAlreadyExists, err)
			},
		},
		{
			name: "should return failed to create verification code",
			ctx:  context.Background(),
			dto:  &registerUserDTO,
			setup: func(ctx context.Context, dto *RegisterUserDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(nil, nil)
				mockUserSvc.EXPECT().CreateUser(ctx, &user.CreateUserDTO{
					Email:    dto.Email,
					Password: dto.Password,
				}).Return("", nil)
				mockVerificationSvc.EXPECT().CreateVerificationCode(ctx, dto.Email).Return("", errors.NewInternal("Failed to create verification code"))
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("Failed to create verification code"), err)
			},
		},
		{
			name: "should return failed to send email",
			ctx:  context.Background(),
			dto:  &registerUserDTO,
			setup: func(ctx context.Context, dto *RegisterUserDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(nil, nil)
				mockUserSvc.EXPECT().CreateUser(ctx, &user.CreateUserDTO{
					Email:    dto.Email,
					Password: dto.Password,
				}).Return("", nil)
				mockVerificationSvc.EXPECT().CreateVerificationCode(ctx, dto.Email).Return(code, nil)
				mockNotificationSvc.EXPECT().SendEmail(ctx, &testEmailData).Return(errors.NewInternal("Failed to send email to user"))
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("Failed to send email to user"), err)
			},
		},
		{
			name: "should return ok",
			ctx:  context.Background(),
			dto:  &registerUserDTO,
			setup: func(ctx context.Context, dto *RegisterUserDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(nil, nil)
				mockUserSvc.EXPECT().CreateUser(ctx, &user.CreateUserDTO{
					Email:    dto.Email,
					Password: dto.Password,
				}).Return("", nil)
				mockVerificationSvc.EXPECT().CreateVerificationCode(ctx, dto.Email).Return(code, nil)
				mockNotificationSvc.EXPECT().SendEmail(ctx, &testEmailData).Return(nil)
			},
			expect: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.dto)
			err := service.RegisterUser(tc.ctx, tc.dto)
			tc.expect(t, err)
		})
	}
}

func TestRegistrationSvc_VerifyUser(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	log := logrus.New()

	mockUserSvc := mock_user.NewMockService(controller)
	mockVerificationSvc := mock_verification.NewMockService(controller)

	deps := &ServiceDeps{
		UserService:         mockUserSvc,
		NotificatorService:  mock_notificator.NewMockService(controller),
		VerificationService: mockVerificationSvc,
		TwoFAService:        mock_twofa.NewMockService(controller),
		JWTService:          mock_jwt.NewMockService(controller),
		CredentialsService:  mock_credentials.NewMockService(controller),
	}

	// Test Data
	emailSender := "example@example.com"
	userEmail := "user@example.com"
	code := "ASDDSA"

	service, _ := NewRegistrationService(log, emailSender, deps)

	// Test DTO
	var verifyUserDTO VerifyUserDTO
	verifyUserDTO.Code = code
	verifyUserDTO.Email = userEmail

	// Test Cred
	secretKey := "secret"
	var testCred credentials.Credentials
	testCred.Password = "==WvZitmZDgzSHgAWvKs"
	testCred.SecretOTP = &secretKey

	// Test wallet
	var testWallet []*wallet.Wallet
	testWallet = append(testWallet, &wallet.Wallet{
		Name:     "BTC",
		WalletId: "8ebdfa95-484d-11ec-ba92-38d547b6cf94",
		Address:  "mrgZBqLCicXRGfSjqiSiV39mXgsV3euVZt",
	})

	// Test user
	testUser, _ := user.NewUser(userEmail, &testWallet, &testCred)

	wrongUserDTO := user.MapToDTO(testUser)
	wrongUserDTO.ID = "example"
	notActiveAndVerifiedUserDTO := user.MapToDTO(testUser)

	testUser.SetToVerified()
	verifiedUserDTO := user.MapToDTO(testUser)

	testUser.SetToActive()

	testUserDTO := user.MapToDTO(testUser)

	tests := []struct {
		name   string
		ctx    context.Context
		dto    *VerifyUserDTO
		setup  func(context.Context, *VerifyUserDTO)
		expect func(*testing.T, error)
	}{
		{
			name: "should return invalid code",
			ctx:  context.Background(),
			dto:  &verifyUserDTO,
			setup: func(ctx context.Context, dto *VerifyUserDTO) {
				mockVerificationSvc.EXPECT().CheckVerificationCode(ctx, dto.Email, dto.Code).Return(errors.NewInternal("Invalid code"))
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("Invalid code"), err)
			},
		},
		{
			name: "should return not found user",
			ctx:  context.Background(),
			dto:  &verifyUserDTO,
			setup: func(ctx context.Context, dto *VerifyUserDTO) {
				mockVerificationSvc.EXPECT().CheckVerificationCode(ctx, dto.Email, dto.Code).Return(nil)
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(nil, errors.NewInternal("User not found"))
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("User not found"), err)
			},
		},
		{
			name: "should return user already verify",
			ctx:  context.Background(),
			dto:  &verifyUserDTO,
			setup: func(ctx context.Context, dto *VerifyUserDTO) {
				mockVerificationSvc.EXPECT().CheckVerificationCode(ctx, dto.Email, dto.Code).Return(nil)
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(testUserDTO, nil)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, ErrAlreadyVerify, err)
			},
		},
		{
			name: "should return wrong object id",
			ctx:  context.Background(),
			dto:  &verifyUserDTO,
			setup: func(ctx context.Context, dto *VerifyUserDTO) {
				mockVerificationSvc.EXPECT().CheckVerificationCode(ctx, dto.Email, dto.Code).Return(nil)
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(wrongUserDTO, nil)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("the provided hex string is not a valid ObjectID"), err)
			},
		},
		{
			name: "should return failed to update user",
			ctx:  context.Background(),
			dto:  &verifyUserDTO,
			setup: func(ctx context.Context, dto *VerifyUserDTO) {
				mockVerificationSvc.EXPECT().CheckVerificationCode(ctx, dto.Email, dto.Code).Return(nil)
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(notActiveAndVerifiedUserDTO, nil)
				mockUserSvc.EXPECT().UpdateUser(ctx, gomock.AssignableToTypeOf(verifiedUserDTO)).Return(errors.NewInternal("Failed to update user"))
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("Failed to update user"), err)
			},
		},
		{
			name: "should return ok",
			ctx:  context.Background(),
			dto:  &verifyUserDTO,
			setup: func(ctx context.Context, dto *VerifyUserDTO) {
				mockVerificationSvc.EXPECT().CheckVerificationCode(ctx, dto.Email, dto.Code).Return(nil)
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(notActiveAndVerifiedUserDTO, nil)
				mockUserSvc.EXPECT().UpdateUser(ctx, gomock.AssignableToTypeOf(verifiedUserDTO)).Return(nil)
			},
			expect: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.dto)
			err := service.VerifyUser(tc.ctx, tc.dto)
			tc.expect(t, err)
		})
	}
}
