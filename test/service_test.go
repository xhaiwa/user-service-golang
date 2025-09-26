package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xhaiwa/user-service-golang/internal/repository"
)

func TestHashPassword(t *testing.T) {
	pass := "mypassword"
	hashed, err := repository.HashPassword(pass)
	assert.NoError(t, err)
	assert.NotEqual(t, pass, hashed)
}

func TestCheckPassword(t *testing.T) {
	pass := "mypassword"
	hashed, _ := repository.HashPassword(pass)
	ok := repository.CheckPassword(hashed, pass)
	assert.True(t, ok)
}
