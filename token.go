// Token Generator
package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go" // JWT library
	"time"
)

// Generate JWT token for a user
func generateToken(username string) (string, error) {
	secretKey := []byte("your-secret-key")
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": username, 
		"timestamp": time.Now().Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}


// Test token generation
func token_main() {
	username := "dcifligu"

	token, err := generateToken(username)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return
	}

	fmt.Println("Generated Token:", token)
}