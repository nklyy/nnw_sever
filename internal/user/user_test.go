package user

import (
	"nnw_s/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		pass      string
		secretOTP string
		salt      int
		expect    func(*testing.T, *User, error)
	}{
		{
			name:      "should return new user",
			email:     "some@mail.com",
			pass:      "somepass",
			secretOTP: "123456",
			salt:      0,
			expect: func(t *testing.T, u *User, err error) {
				assert.NotNil(t, u)
				assert.Nil(t, err)
			},
		},
		{
			name:      "should return ErrInvalidEmail",
			email:     "",
			pass:      "somepass",
			secretOTP: "123456",
			salt:      0,
			expect: func(t *testing.T, u *User, err error) {
				assert.Nil(t, u)
				assert.Equal(t, err, errors.WithMessage(ErrInvalidEmail, "should be not empty"))
			},
		},
		{
			name:      "should return ErrInvalidPassword",
			email:     "some@mail.com",
			pass:      "",
			secretOTP: "123456",
			salt:      0,
			expect: func(t *testing.T, u *User, err error) {
				assert.Nil(t, u)
				assert.Equal(t, err, errors.WithMessage(ErrInvalidPassword, "should be not empty"))
			},
		},
		{
			name:      "should return ErrInvalidPassword due to invalid password salt",
			email:     "some@mail.com",
			pass:      "somepass",
			secretOTP: "123456",
			salt:      9999,
			expect: func(t *testing.T, u *User, err error) {
				assert.Nil(t, u)
				assert.NotNil(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := NewUser(tc.email, tc.pass, tc.secretOTP, tc.salt)
			tc.expect(t, u, err)
		})
	}
}
