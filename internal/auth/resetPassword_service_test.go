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

func TestNewResetPasswordService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	emailSender := "example@example.com"

	tests := []struct {
		name        string
		log         *logrus.Logger
		deps        *ServiceDeps
		emailSender string
		expect      func(*testing.T, ResetPasswordService, error)
	}{
		{
			name: "should return reset password service",
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
			expect: func(t *testing.T, service ResetPasswordService, err error) {
				assert.NotNil(t, service)
				assert.Nil(t, err)
			},
		},
		{
			name:        "should return invalid service dependencies",
			log:         logrus.New(),
			deps:        nil,
			emailSender: emailSender,
			expect: func(t *testing.T, service ResetPasswordService, err error) {
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
			expect: func(t *testing.T, service ResetPasswordService, err error) {
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
			expect: func(t *testing.T, service ResetPasswordService, err error) {
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
			expect: func(t *testing.T, service ResetPasswordService, err error) {
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
			expect: func(t *testing.T, service ResetPasswordService, err error) {
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
			expect: func(t *testing.T, service ResetPasswordService, err error) {
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
			expect: func(t *testing.T, service ResetPasswordService, err error) {
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
			expect: func(t *testing.T, service ResetPasswordService, err error) {
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
			expect: func(t *testing.T, service ResetPasswordService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid sender's email")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := NewResetPasswordService(tc.log, tc.emailSender, tc.deps)
			tc.expect(t, svc, err)
		})
	}
}

func TestResetPasswordSvc_ResetPassword(t *testing.T) {
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
	code := "ASDDSA"

	service, _ := NewResetPasswordService(log, emailSender, deps)

	var resetPasswordDTO ResetPasswordDTO
	resetPasswordDTO.Email = userEmail

	// Test Email Data
	emailData := notificator.Email{
		Subject:   emailResetPasswordSubject,
		Recipient: userEmail,
		Sender:    emailSender,
		Template:  emailResetPasswordTemplateName,
		Data: map[string]interface{}{
			"topic":   emailResetPasswordTopic,
			"message": emailResetPasswordMessage,
			"code":    code,
		},
	}

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

	notActiveUser := user.MapToDTO(testUser)

	testUser.SetToActive()
	testUser.SetToVerified()
	testUserDTO := user.MapToDTO(testUser)

	wrongUserDTO := user.MapToDTO(testUser)
	wrongUserDTO.ID = "example"

	tests := []struct {
		name   string
		ctx    context.Context
		dto    *ResetPasswordDTO
		setup  func(context.Context, *ResetPasswordDTO)
		expect func(*testing.T, error)
	}{
		{
			name: "should return permission_denied",
			ctx:  context.Background(),
			dto:  &resetPasswordDTO,
			setup: func(ctx context.Context, dto *ResetPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(nil, ErrPermissionDenied)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.WithMessage(ErrPermissionDenied, "code: 403; status: permission_denied"), err)
			},
		},
		{
			name: "should return wrong object id",
			ctx:  context.Background(),
			dto:  &resetPasswordDTO,
			setup: func(ctx context.Context, dto *ResetPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(wrongUserDTO, nil)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("the provided hex string is not a valid ObjectID"), err)
			},
		},
		{
			name: "should return user_does_not_verify",
			ctx:  context.Background(),
			dto:  &resetPasswordDTO,
			setup: func(ctx context.Context, dto *ResetPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(notActiveUser, nil)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.WithMessage(user.ErrUserDoesNotVerify, ""), err)
			},
		},
		{
			name: "should return failed to create reset password code",
			ctx:  context.Background(),
			dto:  &resetPasswordDTO,
			setup: func(ctx context.Context, dto *ResetPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(testUserDTO, nil)
				mockVerificationSvc.EXPECT().CreateResetPasswordCode(ctx, dto.Email).Return("", errors.NewInternal("Failed to create reset password code"))
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("Failed to create reset password code"), err)
			},
		},
		{
			name: "should return failed to send email",
			ctx:  context.Background(),
			dto:  &resetPasswordDTO,
			setup: func(ctx context.Context, dto *ResetPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(testUserDTO, nil)
				mockVerificationSvc.EXPECT().CreateResetPasswordCode(ctx, dto.Email).Return(code, nil)
				mockNotificationSvc.EXPECT().SendEmail(ctx, &emailData).Return(errors.NewInternal("Failed to send email"))
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("Failed to send email"), err)
			},
		},
		{
			name: "should return ok",
			ctx:  context.Background(),
			dto:  &resetPasswordDTO,
			setup: func(ctx context.Context, dto *ResetPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(testUserDTO, nil)
				mockVerificationSvc.EXPECT().CreateResetPasswordCode(ctx, dto.Email).Return(code, nil)
				mockNotificationSvc.EXPECT().SendEmail(ctx, &emailData).Return(nil)
			},
			expect: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.dto)
			err := service.ResetPassword(tc.ctx, tc.dto)
			tc.expect(t, err)
		})
	}
}

func TestResetPasswordSvc_ResendResetPasswordEmail(t *testing.T) {
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
	code := "ASDDSA"

	service, _ := NewResetPasswordService(log, emailSender, deps)

	var resendPasswordDTO ResendResetPasswordDTO
	resendPasswordDTO.Email = userEmail

	// Test Email Data
	emailData := notificator.Email{
		Subject:   emailResetPasswordSubject,
		Recipient: userEmail,
		Sender:    emailSender,
		Template:  emailResetPasswordTemplateName,
		Data: map[string]interface{}{
			"topic":   emailResetPasswordTopic,
			"message": emailResetPasswordMessage,
			"code":    code,
		},
	}

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

	notActiveUser := user.MapToDTO(testUser)

	testUser.SetToActive()
	testUser.SetToVerified()
	testUserDTO := user.MapToDTO(testUser)

	wrongUserDTO := user.MapToDTO(testUser)
	wrongUserDTO.ID = "example"

	tests := []struct {
		name   string
		ctx    context.Context
		dto    *ResendResetPasswordDTO
		setup  func(context.Context, *ResendResetPasswordDTO)
		expect func(*testing.T, error)
	}{
		{
			name: "should return permission_denied",
			ctx:  context.Background(),
			dto:  &resendPasswordDTO,
			setup: func(ctx context.Context, dto *ResendResetPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(nil, ErrPermissionDenied)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.WithMessage(ErrPermissionDenied, "code: 403; status: permission_denied"), err)
			},
		},
		{
			name: "should return wrong object id",
			ctx:  context.Background(),
			dto:  &resendPasswordDTO,
			setup: func(ctx context.Context, dto *ResendResetPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(wrongUserDTO, nil)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("the provided hex string is not a valid ObjectID"), err)
			},
		},
		{
			name: "should return user_does_not_verify",
			ctx:  context.Background(),
			dto:  &resendPasswordDTO,
			setup: func(ctx context.Context, dto *ResendResetPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(notActiveUser, nil)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.WithMessage(user.ErrUserDoesNotVerify, ""), err)
			},
		},
		{
			name: "should return failed to create reset password code",
			ctx:  context.Background(),
			dto:  &resendPasswordDTO,
			setup: func(ctx context.Context, dto *ResendResetPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(testUserDTO, nil)
				mockVerificationSvc.EXPECT().CreateResetPasswordCode(ctx, dto.Email).Return("", errors.NewInternal("Failed to create reset password code"))
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("Failed to create reset password code"), err)
			},
		},
		{
			name: "should return failed to send email",
			ctx:  context.Background(),
			dto:  &resendPasswordDTO,
			setup: func(ctx context.Context, dto *ResendResetPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(testUserDTO, nil)
				mockVerificationSvc.EXPECT().CreateResetPasswordCode(ctx, dto.Email).Return(code, nil)
				mockNotificationSvc.EXPECT().SendEmail(ctx, &emailData).Return(errors.NewInternal("Failed to send email"))
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("Failed to send email"), err)
			},
		},
		{
			name: "should return ok",
			ctx:  context.Background(),
			dto:  &resendPasswordDTO,
			setup: func(ctx context.Context, dto *ResendResetPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(testUserDTO, nil)
				mockVerificationSvc.EXPECT().CreateResetPasswordCode(ctx, dto.Email).Return(code, nil)
				mockNotificationSvc.EXPECT().SendEmail(ctx, &emailData).Return(nil)
			},
			expect: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.dto)
			err := service.ResendResetPasswordEmail(tc.ctx, tc.dto)
			tc.expect(t, err)
		})
	}
}

func TestResetPasswordSvc_ResetPasswordCode(t *testing.T) {
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

	service, _ := NewResetPasswordService(log, emailSender, deps)

	// Test DTO
	var resetPasswordCodeDTO ResetPasswordCodedDTO
	resetPasswordCodeDTO.Code = "ASDDSA"
	resetPasswordCodeDTO.Email = userEmail

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

	notActiveUser := user.MapToDTO(testUser)

	testUser.SetToActive()
	testUser.SetToVerified()
	testUserDTO := user.MapToDTO(testUser)

	wrongUserDTO := user.MapToDTO(testUser)
	wrongUserDTO.ID = "example"

	tests := []struct {
		name   string
		ctx    context.Context
		dto    *ResetPasswordCodedDTO
		setup  func(context.Context, *ResetPasswordCodedDTO)
		expect func(*testing.T, error)
	}{
		{
			name: "should return permission_denied",
			ctx:  context.Background(),
			dto:  &resetPasswordCodeDTO,
			setup: func(ctx context.Context, dto *ResetPasswordCodedDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(nil, ErrPermissionDenied)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.WithMessage(ErrPermissionDenied, "code: 403; status: permission_denied"), err)
			},
		},
		{
			name: "should return wrong object id",
			ctx:  context.Background(),
			dto:  &resetPasswordCodeDTO,
			setup: func(ctx context.Context, dto *ResetPasswordCodedDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(wrongUserDTO, nil)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("the provided hex string is not a valid ObjectID"), err)
			},
		},
		{
			name: "should return permission_denied not active user",
			ctx:  context.Background(),
			dto:  &resetPasswordCodeDTO,
			setup: func(ctx context.Context, dto *ResetPasswordCodedDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(notActiveUser, nil)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.WithMessage(ErrPermissionDenied, ""), err)
			},
		},
		{
			name: "should return check reset password code",
			ctx:  context.Background(),
			dto:  &resetPasswordCodeDTO,
			setup: func(ctx context.Context, dto *ResetPasswordCodedDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(testUserDTO, nil)
				mockVerificationSvc.EXPECT().CheckResetPasswordCode(ctx, dto.Email, dto.Code).Return(errors.NewInternal("Failed to check reset password code"))
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("Failed to check reset password code"), err)
			},
		},
		{
			name: "should return ok",
			ctx:  context.Background(),
			dto:  &resetPasswordCodeDTO,
			setup: func(ctx context.Context, dto *ResetPasswordCodedDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(testUserDTO, nil)
				mockVerificationSvc.EXPECT().CheckResetPasswordCode(ctx, dto.Email, dto.Code).Return(nil)
			},
			expect: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.dto)
			err := service.ResetPasswordCode(tc.ctx, tc.dto)
			tc.expect(t, err)
		})
	}
}

func TestResetPasswordSvc_SetupNewPassword(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	log := logrus.New()

	mockUserSvc := mock_user.NewMockService(controller)
	mockCredentialsSvc := mock_credentials.NewMockService(controller)

	deps := &ServiceDeps{
		UserService:         mockUserSvc,
		NotificatorService:  mock_notificator.NewMockService(controller),
		VerificationService: mock_verification.NewMockService(controller),
		TwoFAService:        mock_twofa.NewMockService(controller),
		JWTService:          mock_jwt.NewMockService(controller),
		CredentialsService:  mockCredentialsSvc,
	}

	// Test Data
	emailSender := "example@example.com"
	userEmail := "user@example.com"
	userPassword := "==WvZitmZDgzSHgAWvKs"

	service, _ := NewResetPasswordService(log, emailSender, deps)

	// Test DTO
	var setupNewPasswordDTO SetupNewPasswordDTO
	setupNewPasswordDTO.Email = userEmail
	setupNewPasswordDTO.Password = userPassword

	// Test Cred
	secretKey := "secret"
	var testCred credentials.Credentials
	testCred.Password = "==WvZitmZDgzSHgAWvKs"
	testCred.SecretOTP = &secretKey

	testCredDTO := credentials.MapToDTO(&testCred)

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

	tests := []struct {
		name   string
		ctx    context.Context
		dto    *SetupNewPasswordDTO
		setup  func(context.Context, *SetupNewPasswordDTO)
		expect func(*testing.T, error)
	}{
		{
			name: "should return permission_denied",
			ctx:  context.Background(),
			dto:  &setupNewPasswordDTO,
			setup: func(ctx context.Context, dto *SetupNewPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(nil, ErrPermissionDenied)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.WithMessage(ErrPermissionDenied, "code: 403; status: permission_denied"), err)
			},
		},
		{
			name: "should return wrong object id",
			ctx:  context.Background(),
			dto:  &setupNewPasswordDTO,
			setup: func(ctx context.Context, dto *SetupNewPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(wrongUserDTO, nil)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("the provided hex string is not a valid ObjectID"), err)
			},
		},
		{
			name: "should return permission_denied not active user",
			ctx:  context.Background(),
			dto:  &setupNewPasswordDTO,
			setup: func(ctx context.Context, dto *SetupNewPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(notActiveUser, nil)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.WithMessage(ErrPermissionDenied, ""), err)
			},
		},
		{
			name: "should return failed to create user credentials",
			ctx:  context.Background(),
			dto:  &setupNewPasswordDTO,
			setup: func(ctx context.Context, dto *SetupNewPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(testUserDTO, nil)
				mockCredentialsSvc.EXPECT().CreateCredentials(ctx, dto.Password, testCred.SecretOTP).Return(nil, errors.NewInternal("Failed to create user credentials"))
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("Failed to create user credentials"), err)
			},
		},
		{
			name: "should return failed to update user",
			ctx:  context.Background(),
			dto:  &setupNewPasswordDTO,
			setup: func(ctx context.Context, dto *SetupNewPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(testUserDTO, nil)
				mockCredentialsSvc.EXPECT().CreateCredentials(ctx, dto.Password, testCred.SecretOTP).Return(testCredDTO, nil)
				mockUserSvc.EXPECT().UpdateUser(ctx, testUserDTO).Return(errors.NewInternal("Failed to update user"))
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("Failed to update user"), err)
			},
		},
		{
			name: "should return ok",
			ctx:  context.Background(),
			dto:  &setupNewPasswordDTO,
			setup: func(ctx context.Context, dto *SetupNewPasswordDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(testUserDTO, nil)
				mockCredentialsSvc.EXPECT().CreateCredentials(ctx, dto.Password, testCred.SecretOTP).Return(testCredDTO, nil)
				mockUserSvc.EXPECT().UpdateUser(ctx, testUserDTO).Return(nil)
			},
			expect: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.dto)
			err := service.SetupNewPassword(tc.ctx, tc.dto)
			tc.expect(t, err)
		})
	}
}
