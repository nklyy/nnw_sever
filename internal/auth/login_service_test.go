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
	mock_notificator "nnw_s/pkg/notificator/mocks"
	"testing"
)

func TestNewLoginService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name   string
		log    *logrus.Logger
		deps   *ServiceDeps
		expect func(*testing.T, LoginService, error)
	}{
		{
			name: "should return login service",
			log:  logrus.New(),
			deps: &ServiceDeps{
				UserService:         mock_user.NewMockService(controller),
				NotificatorService:  mock_notificator.NewMockService(controller),
				VerificationService: mock_verification.NewMockService(controller),
				TwoFAService:        mock_twofa.NewMockService(controller),
				JWTService:          mock_jwt.NewMockService(controller),
				CredentialsService:  mock_credentials.NewMockService(controller),
			},
			expect: func(t *testing.T, service LoginService, err error) {
				assert.NotNil(t, service)
				assert.Nil(t, err)
			},
		},
		{
			name: "should return 'invalid service dependencies'",
			log:  logrus.New(),
			deps: nil,
			expect: func(t *testing.T, service LoginService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid service dependencies")
			},
		},
		{
			name: "should return 'invalid user service'",
			log:  logrus.New(),
			deps: &ServiceDeps{
				UserService:         nil,
				NotificatorService:  mock_notificator.NewMockService(controller),
				VerificationService: mock_verification.NewMockService(controller),
				TwoFAService:        mock_twofa.NewMockService(controller),
				JWTService:          mock_jwt.NewMockService(controller),
				CredentialsService:  mock_credentials.NewMockService(controller),
			},
			expect: func(t *testing.T, service LoginService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid user service")
			},
		},
		{
			name: "should return 'invalid notification service'",
			log:  logrus.New(),
			deps: &ServiceDeps{
				UserService:         mock_user.NewMockService(controller),
				NotificatorService:  nil,
				VerificationService: mock_verification.NewMockService(controller),
				TwoFAService:        mock_twofa.NewMockService(controller),
				JWTService:          mock_jwt.NewMockService(controller),
				CredentialsService:  mock_credentials.NewMockService(controller),
			},
			expect: func(t *testing.T, service LoginService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid notification service")
			},
		},
		{
			name: "should return 'invalid verification service'",
			log:  logrus.New(),
			deps: &ServiceDeps{
				UserService:         mock_user.NewMockService(controller),
				NotificatorService:  mock_notificator.NewMockService(controller),
				VerificationService: nil,
				TwoFAService:        mock_twofa.NewMockService(controller),
				JWTService:          mock_jwt.NewMockService(controller),
				CredentialsService:  mock_credentials.NewMockService(controller),
			},
			expect: func(t *testing.T, service LoginService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid verification service")
			},
		},
		{
			name: "should return 'invalid TwoFA service'",
			log:  logrus.New(),
			deps: &ServiceDeps{
				UserService:         mock_user.NewMockService(controller),
				NotificatorService:  mock_notificator.NewMockService(controller),
				VerificationService: mock_verification.NewMockService(controller),
				TwoFAService:        nil,
				JWTService:          mock_jwt.NewMockService(controller),
				CredentialsService:  mock_credentials.NewMockService(controller),
			},
			expect: func(t *testing.T, service LoginService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid TwoFA service")
			},
		},
		{
			name: "should return 'invalid JWT service'",
			log:  logrus.New(),
			deps: &ServiceDeps{
				UserService:         mock_user.NewMockService(controller),
				NotificatorService:  mock_notificator.NewMockService(controller),
				VerificationService: mock_verification.NewMockService(controller),
				TwoFAService:        mock_twofa.NewMockService(controller),
				JWTService:          nil,
				CredentialsService:  mock_credentials.NewMockService(controller),
			},
			expect: func(t *testing.T, service LoginService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid JWT service")
			},
		},
		{
			name: "should return 'invalid credentials service'",
			log:  logrus.New(),
			deps: &ServiceDeps{
				UserService:         mock_user.NewMockService(controller),
				NotificatorService:  mock_notificator.NewMockService(controller),
				VerificationService: mock_verification.NewMockService(controller),
				TwoFAService:        mock_twofa.NewMockService(controller),
				JWTService:          mock_jwt.NewMockService(controller),
				CredentialsService:  nil,
			},
			expect: func(t *testing.T, service LoginService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid credentials service")
			},
		},
		{
			name: "should return 'invalid logger'",
			log:  nil,
			deps: &ServiceDeps{
				UserService:         mock_user.NewMockService(controller),
				NotificatorService:  mock_notificator.NewMockService(controller),
				VerificationService: mock_verification.NewMockService(controller),
				TwoFAService:        mock_twofa.NewMockService(controller),
				JWTService:          mock_jwt.NewMockService(controller),
				CredentialsService:  mock_credentials.NewMockService(controller),
			},
			expect: func(t *testing.T, service LoginService, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid logger")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := NewLoginService(tc.log, tc.deps)
			tc.expect(t, svc, err)
		})
	}
}

func TestLoginSvc_Login(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	log := logrus.New()
	deps := &ServiceDeps{
		UserService:         mock_user.NewMockService(controller),
		NotificatorService:  mock_notificator.NewMockService(controller),
		VerificationService: mock_verification.NewMockService(controller),
		TwoFAService:        mock_twofa.NewMockService(controller),
		JWTService:          mock_jwt.NewMockService(controller),
		CredentialsService:  mock_credentials.NewMockService(controller),
	}

	service, _ := NewLoginService(log, deps)
	mockUserSvc := mock_user.NewMockService(controller)
	mockCredSvc := mock_credentials.NewMockService(controller)

	var loginUserDTO LoginDTO
	loginUserDTO.Email = "some@mail.com"
	loginUserDTO.Password = "==WvZitmZDgzSHgAWvKs"

	secretKey := "secret"
	var testCred credentials.Credentials
	testCred.Password = "==WvZitmZDgzSHgAWvKs"
	testCred.SecretOTP = &secretKey
	credDTO := credentials.MapToDTO(&testCred)

	testUser, _ := user.NewUser("some@mail.com", &testCred)
	userDTO := user.MapToDTO(testUser)

	tests := []struct {
		name   string
		ctx    context.Context
		dto    *LoginDTO
		setup  func(context.Context, *LoginDTO)
		expect func(*testing.T, error)
	}{
		{
			name: "should return status ok",
			ctx:  context.Background(),
			dto:  &loginUserDTO,
			setup: func(ctx context.Context, dto *LoginDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(userDTO, nil)
				mockCredSvc.EXPECT().ValidatePassword(ctx, credDTO, dto.Password).Return(nil)
			},
			expect: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.dto)
			err := service.Login(tc.ctx, tc.dto)
			tc.expect(t, err)
		})
	}
}
