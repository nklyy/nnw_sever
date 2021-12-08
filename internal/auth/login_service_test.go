package auth

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"nnw_s/internal/auth/jwt"
	mock_jwt "nnw_s/internal/auth/jwt/mocks"
	"nnw_s/internal/auth/twofa"
	mock_twofa "nnw_s/internal/auth/twofa/mocks"
	mock_verification "nnw_s/internal/auth/verification/mocks"
	"nnw_s/internal/user"
	"nnw_s/internal/user/credentials"
	mock_credentials "nnw_s/internal/user/credentials/mocks"
	mock_user "nnw_s/internal/user/mocks"
	"nnw_s/pkg/errors"
	mock_notificator "nnw_s/pkg/notificator/mocks"
	"nnw_s/pkg/wallet"
	"testing"
	"time"
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
			name: "should return 'invalid user service'",
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

	// Test wallet
	var testWallet []*wallet.Wallet
	testWallet = append(testWallet, &wallet.Wallet{
		Name:     "BTC",
		WalletId: "8ebdfa95-484d-11ec-ba92-38d547b6cf94",
		Address:  "mrgZBqLCicXRGfSjqiSiV39mXgsV3euVZt",
	})

	// Verify user
	testActiveUser, _ := user.NewUser("some@mail.com", &testWallet, &testCred)
	testActiveUser.SetToActive()
	testActiveUser.SetToVerified()
	activeUserDTO := user.MapToDTO(testActiveUser)

	// Disable user
	testDisableUser, _ := user.NewUser("some@mail.com", &testWallet, &testCred)
	disableUserDTO := user.MapToDTO(testDisableUser)

	// User with wrong id
	testWrongUser, _ := user.NewUser("some@mail.com", &testWallet, &testCred)
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
			name: "should wrong object id",
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

func TestLoginSvc_CheckCode(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	log := logrus.New()
	mockUserSvc := mock_user.NewMockService(controller)
	mockTwoFaSvc := mock_twofa.NewMockService(controller)
	mockJwtSvc := mock_jwt.NewMockService(controller)
	deps := &ServiceDeps{
		UserService:         mockUserSvc,
		NotificatorService:  mock_notificator.NewMockService(controller),
		VerificationService: mock_verification.NewMockService(controller),
		TwoFAService:        mockTwoFaSvc,
		JWTService:          mockJwtSvc,
		CredentialsService:  mock_credentials.NewMockService(controller),
	}

	service, _ := NewLoginService(log, deps)

	// Cred
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

	// Verify user
	testActiveUser, _ := user.NewUser("some@mail.com", &testWallet, &testCred)
	testActiveUser.SetToActive()
	testActiveUser.SetToVerified()
	activeUserDTO := user.MapToDTO(testActiveUser)

	// Disable user
	testDisableUser, _ := user.NewUser("some@mail.com", &testWallet, &testCred)
	disableUserDTO := user.MapToDTO(testDisableUser)

	// User with wrong id
	testWrongUser, _ := user.NewUser("some@mail.com", &testWallet, &testCred)
	testWrongUser.SetToActive()
	testWrongUser.SetToVerified()
	wrongUserDTO := user.MapToDTO(testWrongUser)
	wrongUserDTO.ID = "example"

	var loginCodeDTO LoginCodeDTO
	loginCodeDTO.Email = "some@mail.com"
	loginCodeDTO.Code = "241241"

	var testJwtDTO jwt.DTO
	testJwtDTO.ID = "id"
	testJwtDTO.Token = "token"
	testJwtDTO.ExpireAt = time.Now()

	tests := []struct {
		name     string
		ctx      context.Context
		loginDto *LoginCodeDTO
		setup    func(context.Context, *LoginCodeDTO)
		expect   func(*testing.T, *TokenDTO, error)
	}{
		{
			name:     "should return token",
			ctx:      context.Background(),
			loginDto: &loginCodeDTO,
			setup: func(ctx context.Context, loginDto *LoginCodeDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, loginDto.Email).Return(activeUserDTO, nil)
				mockTwoFaSvc.EXPECT().CheckTwoFACode(ctx, loginDto.Code, *testCred.SecretOTP).Return(nil)
				mockJwtSvc.EXPECT().CreateJWT(ctx, loginDto.Email).Return(&testJwtDTO, nil)
			},
			expect: func(t *testing.T, dto *TokenDTO, err error) {
				assert.NotEmpty(t, dto)
				assert.Nil(t, err)
				assert.Equal(t, dto.Token, testJwtDTO.Token)
			},
		},
		{
			name:     "should permission_denied by getUserByEmail",
			ctx:      context.Background(),
			loginDto: &loginCodeDTO,
			setup: func(ctx context.Context, loginDto *LoginCodeDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, loginDto.Email).Return(nil, ErrPermissionDenied)
			},
			expect: func(t *testing.T, dto *TokenDTO, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.WithMessage(ErrPermissionDenied, "code: 403; status: permission_denied"), err)
			},
		},
		{
			name:     "should wrong object id",
			ctx:      context.Background(),
			loginDto: &loginCodeDTO,
			setup: func(ctx context.Context, loginDto *LoginCodeDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, loginDto.Email).Return(wrongUserDTO, nil)
			},
			expect: func(t *testing.T, dto *TokenDTO, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.NewInternal("the provided hex string is not a valid ObjectID"), err)
			},
		},
		{
			name:     "should permission_denied disable user",
			ctx:      context.Background(),
			loginDto: &loginCodeDTO,
			setup: func(ctx context.Context, loginDto *LoginCodeDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, loginDto.Email).Return(disableUserDTO, nil)
			},
			expect: func(t *testing.T, dto *TokenDTO, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, ErrPermissionDenied, err)
			},
		},
		{
			name:     "should invalid twoFa code",
			ctx:      context.Background(),
			loginDto: &loginCodeDTO,
			setup: func(ctx context.Context, loginDto *LoginCodeDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, loginDto.Email).Return(activeUserDTO, nil)
				mockTwoFaSvc.EXPECT().CheckTwoFACode(ctx, loginDto.Code, *testCred.SecretOTP).Return(twofa.ErrInvalidTwoFACode)
			},
			expect: func(t *testing.T, dto *TokenDTO, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, twofa.ErrInvalidTwoFACode, err)
			},
		},
		{
			name:     "should token invalid",
			ctx:      context.Background(),
			loginDto: &loginCodeDTO,
			setup: func(ctx context.Context, loginDto *LoginCodeDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, loginDto.Email).Return(activeUserDTO, nil)
				mockTwoFaSvc.EXPECT().CheckTwoFACode(ctx, loginDto.Code, *testCred.SecretOTP).Return(nil)
				mockJwtSvc.EXPECT().CreateJWT(ctx, loginDto.Email).Return(nil, jwt.ErrTokenDoesNotValid)
			},
			expect: func(t *testing.T, dto *TokenDTO, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.WithMessage(ErrUnauthorized, "code: 401; status: token_invalid"), err)
			},
		},
		{
			name:     "should token expired",
			ctx:      context.Background(),
			loginDto: &loginCodeDTO,
			setup: func(ctx context.Context, loginDto *LoginCodeDTO) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, loginDto.Email).Return(activeUserDTO, nil)
				mockTwoFaSvc.EXPECT().CheckTwoFACode(ctx, loginDto.Code, *testCred.SecretOTP).Return(nil)
				mockJwtSvc.EXPECT().CreateJWT(ctx, loginDto.Email).Return(nil, jwt.ErrTokenHasBeenExpired)
			},
			expect: func(t *testing.T, dto *TokenDTO, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, errors.WithMessage(ErrUnauthorized, "code: 401; status: token_expired"), err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.loginDto)
			tokenDTO, err := service.CheckCode(tc.ctx, tc.loginDto)
			tc.expect(t, tokenDTO, err)
		})
	}
}
