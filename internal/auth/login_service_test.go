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
	mockUserSvc := mock_user.NewMockService(controller)
	mockCredSvc := mock_credentials.NewMockService(controller)
	deps := &ServiceDeps{
		UserService:         mockUserSvc,
		NotificatorService:  mock_notificator.NewMockService(controller),
		VerificationService: mock_verification.NewMockService(controller),
		TwoFAService:        mock_twofa.NewMockService(controller),
		JWTService:          mock_jwt.NewMockService(controller),
		CredentialsService:  mockCredSvc,
	}

	service, _ := NewLoginService(log, deps)

	var loginUserDTO LoginDTO
	loginUserDTO.Email = "some@mail.com"
	loginUserDTO.Password = "==WvZitmZDgzSHgAWvKs"

	// Cred
	secretKey := "secret"
	var testCred credentials.Credentials
	testCred.Password = "==WvZitmZDgzSHgAWvKs"
	testCred.SecretOTP = &secretKey
	credDTO := credentials.MapToDTO(&testCred)

	// Verify user
	testActiveUser, _ := user.NewUser("some@mail.com", &testCred)
	testActiveUser.SetToActive()
	testActiveUser.SetToVerified()
	activeUserDTO := user.MapToDTO(testActiveUser)

	// Disable user
	testDisableUser, _ := user.NewUser("some@mail.com", &testCred)
	disableUserDTO := user.MapToDTO(testDisableUser)

	// User with wrong id
	testWrongUser, _ := user.NewUser("some@mail.com", &testCred)
	testWrongUser.SetToActive()
	testWrongUser.SetToVerified()
	wrongUserDTO := user.MapToDTO(testWrongUser)
	wrongUserDTO.ID = "example"

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
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(activeUserDTO, nil)
				mockCredSvc.EXPECT().ValidatePassword(ctx, credDTO, dto.Password).Return(nil)
			},
			expect: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "should permission_denied by getUserByEmail",
			ctx:  context.Background(),
			dto:  &loginUserDTO,
			setup: func(ctx context.Context, dto *LoginDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(nil, ErrPermissionDenied)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.WithMessage(ErrPermissionDenied, "code: 403; status: permission_denied"), err)
			},
		},
		{
			name: "should permission_denied disable user",
			ctx:  context.Background(),
			dto:  &loginUserDTO,
			setup: func(ctx context.Context, dto *LoginDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(disableUserDTO, nil)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, ErrPermissionDenied, err)
			},
		},
		{
			name: "should invalid_user_password",
			ctx:  context.Background(),
			dto:  &loginUserDTO,
			setup: func(ctx context.Context, dto *LoginDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(activeUserDTO, nil)
				mockCredSvc.EXPECT().ValidatePassword(ctx, credDTO, dto.Password).Return(user.ErrInvalidPassword)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, user.ErrInvalidPassword, err)
			},
		},
		{
			name: "should invalid_user_password",
			ctx:  context.Background(),
			dto:  &loginUserDTO,
			setup: func(ctx context.Context, dto *LoginDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email).Return(wrongUserDTO, nil)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("the provided hex string is not a valid ObjectID"), err)
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
