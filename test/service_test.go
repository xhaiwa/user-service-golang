package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xhaiwa/user-service-golang/internal/service"
)

func TestHashPassword(t *testing.T) {
	pass := "mypassword"
	hashed, err := service.HashPassword(pass)
	assert.NoError(t, err)
	assert.NotEqual(t, pass, hashed)
}

func TestCheckPassword(t *testing.T) {
	pass := "mypassword"
	hashed, _ := service.HashPassword(pass)
	ok := service.CheckPassword(hashed, pass)
	assert.True(t, ok)
}
