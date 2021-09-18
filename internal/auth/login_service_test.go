package auth

import (
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	mock_jwt "nnw_s/internal/auth/jwt/mocks"
	mock_twofa "nnw_s/internal/auth/twofa/mocks"
	mock_verification "nnw_s/internal/auth/verification/mocks"
	mock_credentials "nnw_s/internal/user/credentials/mocks"
	mock_user "nnw_s/internal/user/mocks"
	mock_notificator "nnw_s/pkg/notificator/mocks"
	"testing"
)

func TestNewLoginService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	testDep := ServiceDeps{
		UserService:         mock_user.NewMockService(controller),
		NotificatorService:  mock_notificator.NewMockService(controller),
		VerificationService: mock_verification.NewMockService(controller),
		TwoFAService:        mock_twofa.NewMockService(controller),
		JWTService:          mock_jwt.NewMockService(controller),
		CredentialsService:  mock_credentials.NewMockService(controller),
	}

	tests := []struct {
		name   string
		log    *logrus.Logger
		deps   *ServiceDeps
		expect func(*testing.T, LoginService, error)
	}{
		{
			name: "should return login service",
			log:  logrus.New(),
			deps: &testDep,
			expect: func(t *testing.T, service LoginService, err error) {
				assert.NotNil(t, service)
				assert.Nil(t, err)
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
