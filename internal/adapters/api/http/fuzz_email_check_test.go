package http

import (
	"net/mail"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FuzzValidateEmail(f *testing.F) {
	testcases := []string{"test@example.com", "invalid-email", "another@test.co", "ars-saz@ya.ru", ""}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, email string) {
		valid := validateEmail(email)
		assert.Condition(t, func() (success bool) {
			_, err := mail.ParseAddress(email)
			return valid == (err == nil)
		}, "Email validation mismatch")
	})
}
