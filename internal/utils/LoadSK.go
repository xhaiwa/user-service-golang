package utils

import (
	"log"
	"os"
)

var JWTSecret string

func LoadSK() {
	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
}
