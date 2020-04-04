package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	cases := []struct {
		expected bool
		email    string
	}{
		{true, "test@test.com"},
		{true, "1@1.com"},
		{false, ""},
		{true, "-@gmail.com"},
		{true, "donald.duck@gmail.com"},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("case %d: '%s'", i, tt.email), func(t *testing.T) {
			assert.Equal(t, tt.expected, isEmailValid(tt.email))
		})
	}
}
