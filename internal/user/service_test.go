package user_test

import (
	"context"
	"nnw_s/internal/user"
	mock_user "nnw_s/internal/user/mocks"
	"nnw_s/pkg/errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	testOpts := &user.ServiceOptions{
		Log:          logrus.New(),
		Shift:        5,
		PasswordSalt: 5,
	}

	tests := []struct {
		name   string
		repo   user.Repository
		opts   *user.ServiceOptions
		expect func(*testing.T, user.Service, error)
	}{
		{
			name: "should return user service",
			repo: mock_user.NewMockRepository(controller),
			opts: testOpts,
			expect: func(t *testing.T, s user.Service, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name: "should return 'invalid repo' error",
			repo: nil,
			opts: testOpts,
			expect: func(t *testing.T, s user.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid repo")
			},
		},
		{
			name: "should return 'invalid service options' error",
			repo: mock_user.NewMockRepository(controller),
			opts: nil,
			expect: func(t *testing.T, s user.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid service options")
			},
		},
		{
			name: "should return 'invalid logger' error",
			repo: mock_user.NewMockRepository(controller),
			opts: &user.ServiceOptions{
				Shift:        testOpts.Shift,
				PasswordSalt: testOpts.PasswordSalt,
			},
			expect: func(t *testing.T, s user.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid logger")
			},
		},
		{
			name: "should return 'invalid password salt' error",
			repo: mock_user.NewMockRepository(controller),
			opts: &user.ServiceOptions{
				Log:   testOpts.Log,
				Shift: testOpts.Shift,
			},
			expect: func(t *testing.T, s user.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid password salt")
			},
		},
		{
			name: "should return 'invalid shift' error",
			repo: mock_user.NewMockRepository(controller),
			opts: &user.ServiceOptions{
				Log:          testOpts.Log,
				PasswordSalt: testOpts.PasswordSalt,
			},
			expect: func(t *testing.T, s user.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "code: 500; status: internal_error; message: invalid shift")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := user.NewService(tc.repo, tc.opts)
			tc.expect(t, svc, err)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_user.NewMockRepository(controller)
	testOpts := &user.ServiceOptions{
		Log:          logrus.New(),
		Shift:        5,
		PasswordSalt: 5,
	}

	service, _ := user.NewService(mockRepo, testOpts)
	testUser, _ := user.NewUser("some@mail.com", "Password12345", "secret")

	tests := []struct {
		name   string
		ctx    context.Context
		userID string
		setup  func(context.Context, string)
		expect func(*testing.T, *user.User, error)
	}{
		{
			name:   "should return test user",
			ctx:    context.Background(),
			userID: testUser.ID.Hex(),
			setup: func(ctx context.Context, userID string) {
				mockRepo.EXPECT().GetUserByID(ctx, userID).Return(testUser, nil)
			},
			expect: func(t *testing.T, u *user.User, err error) {
				assert.NotNil(t, u)
				assert.Nil(t, err)
				assert.Equal(t, testUser, u)
			},
		},
		{
			name:   "should return 'not found' error",
			ctx:    context.Background(),
			userID: "not_existent_id",
			setup: func(ctx context.Context, userID string) {
				mockRepo.EXPECT().GetUserByID(ctx, userID).Return(nil, user.ErrNotFound)
			},
			expect: func(t *testing.T, u *user.User, err error) {
				assert.Nil(t, u)
				assert.Equal(t, user.ErrNotFound, err)
			},
		},
		{
			name:   "should return 'internal error' error",
			ctx:    context.Background(),
			userID: testUser.ID.Hex(),
			setup: func(ctx context.Context, userID string) {
				mockRepo.EXPECT().GetUserByID(ctx, userID).Return(nil, errors.NewInternal("internal error"))
			},
			expect: func(t *testing.T, u *user.User, err error) {
				assert.Nil(t, u)
				assert.Equal(t, errors.NewInternal("internal error"), err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.userID)
			u, err := service.GetUserByID(tc.ctx, tc.userID)
			tc.expect(t, u, err)
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_user.NewMockRepository(controller)
	testOpts := &user.ServiceOptions{
		Log:          logrus.New(),
		Shift:        5,
		PasswordSalt: 5,
	}

	service, _ := user.NewService(mockRepo, testOpts)
	testUser, _ := user.NewUser("some@mail.com", "Password12345", "secret")

	tests := []struct {
		name   string
		ctx    context.Context
		email  string
		setup  func(context.Context, string)
		expect func(*testing.T, *user.User, error)
	}{
		{
			name:  "should return test user",
			ctx:   context.Background(),
			email: testUser.Email,
			setup: func(ctx context.Context, email string) {
				mockRepo.EXPECT().GetUserByEmail(ctx, email).Return(testUser, nil)
			},
			expect: func(t *testing.T, u *user.User, err error) {
				assert.NotNil(t, u)
				assert.Nil(t, err)
				assert.Equal(t, testUser, u)
			},
		},
		{
			name:  "should return 'not found' error",
			ctx:   context.Background(),
			email: "not_existent_email",
			setup: func(ctx context.Context, email string) {
				mockRepo.EXPECT().GetUserByEmail(ctx, email).Return(nil, user.ErrNotFound)
			},
			expect: func(t *testing.T, u *user.User, err error) {
				assert.Nil(t, u)
				assert.Equal(t, user.ErrNotFound, err)
			},
		},
		{
			name:  "should return 'internal error' error",
			ctx:   context.Background(),
			email: testUser.Email,
			setup: func(ctx context.Context, email string) {
				mockRepo.EXPECT().GetUserByEmail(ctx, email).Return(nil, errors.NewInternal("internal error"))
			},
			expect: func(t *testing.T, u *user.User, err error) {
				assert.Nil(t, u)
				assert.Equal(t, errors.NewInternal("internal error"), err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.email)
			u, err := service.GetUserByEmail(tc.ctx, tc.email)
			tc.expect(t, u, err)
		})
	}
}

func TestCreateUser(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_user.NewMockRepository(controller)
	testOpts := &user.ServiceOptions{
		Log:          logrus.New(),
		Shift:        15,
		PasswordSalt: 5,
	}

	service, _ := user.NewService(mockRepo, testOpts)
	testUser, _ := user.NewUser("some@mail.com", "$2a$05$gNX0NIYmmC/1rPGEclZrVeBR9DZmt4l2ydNn8tM5XNfFJxlL8ObXq", "secret")
	encodedPass := "==WvZitmZDgzSHgAWvKs"

	tests := []struct {
		name   string
		ctx    context.Context
		dto    *user.CreateUserDTO
		setup  func(context.Context, *user.CreateUserDTO)
		expect func(*testing.T, string, error)
	}{
		{
			name: "should return id of recently created user",
			ctx:  context.Background(),
			dto: &user.CreateUserDTO{
				Email:     testUser.Email,
				Password:  encodedPass,
				SecretOTP: testUser.SecretOTP,
			},
			setup: func(ctx context.Context, dto *user.CreateUserDTO) {
				mockRepo.EXPECT().SaveUser(ctx, gomock.Any()).Return(testUser.ID.Hex(), nil)
			},
			expect: func(t *testing.T, id string, err error) {
				assert.NotEmpty(t, id)
				assert.Nil(t, err)
				assert.Equal(t, id, testUser.ID.Hex())
			},
		},
		{
			name: "should return decode error",
			ctx:  context.Background(),
			dto: &user.CreateUserDTO{
				Email:     testUser.Email,
				Password:  testUser.Password,
				SecretOTP: testUser.SecretOTP,
			},
			setup: func(ctx context.Context, dto *user.CreateUserDTO) {},
			expect: func(t *testing.T, id string, err error) {
				assert.Empty(t, id)
				assert.NotNil(t, err)
			},
		},
		{
			name: "should return error while saving user is db",
			ctx:  context.Background(),
			dto: &user.CreateUserDTO{
				Email:     testUser.Email,
				Password:  encodedPass,
				SecretOTP: testUser.SecretOTP,
			},
			setup: func(ctx context.Context, dto *user.CreateUserDTO) {
				mockRepo.EXPECT().SaveUser(ctx, gomock.Any()).Return("", user.ErrAlreadyExists)
			},
			expect: func(t *testing.T, id string, err error) {
				assert.Empty(t, id)
				assert.Equal(t, user.ErrAlreadyExists, err)
			},
		},
		{
			name: "should return error while creating a new user",
			ctx:  context.Background(),
			dto: &user.CreateUserDTO{
				Email:     "",
				Password:  encodedPass,
				SecretOTP: testUser.SecretOTP,
			},
			setup: func(ctx context.Context, dto *user.CreateUserDTO) {},
			expect: func(t *testing.T, id string, err error) {
				assert.Empty(t, id)
				assert.Equal(t, errors.WithMessage(user.ErrInvalidEmail, "should be not empty"), err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.dto)
			id, err := service.CreateUser(tc.ctx, tc.dto)
			tc.expect(t, id, err)
		})
	}
}
