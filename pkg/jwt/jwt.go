package jwt

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// JWTToken is a struct that holds the private and public keys
type JWTToken struct {
	privateKey []byte
	publicKey  []byte
}

var JWT JWTToken

func InitJWTToken(privateKeyPath, publicKeyPath string) {
	// Read private key from file
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	// Read public key from file
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	JWT = JWTToken{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// SignToken signs a token using a private key
func SignToken(body map[string]interface{}) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(JWT.privateKey)
	if err != nil {
		return "", fmt.Errorf("generate token parse key error: %w", err)
	}
	token := jwt.New(jwt.SigningMethodRS256)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	for k, v := range body {
		claims[k] = v
	}
	// Generate encoded token and send it as response.
	return token.SignedString(key)
}

// VerifyToken verifies a token using a public key
func VerifyToken(token string) (map[string]interface{}, error) {
	// Parse the private key
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(JWT.publicKey)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %v", err)
	}

	tok, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate invalid")
	}

	return claims, nil
}
