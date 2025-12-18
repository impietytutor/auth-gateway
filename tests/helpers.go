package authgateway

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// GenerateRSAKeyPair generates a new RSA key pair.
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, nil
}

// PEMEncrypt encrypts a given plaintext using a provided public key.
func PEMEncrypt(publicKey *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	if publicKey == nil {
		return nil, errors.New("public key cannot be nil")
	}
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt plaintext: %w", err)
	}
	return encrypted, nil
}

// PEMDecrypt decrypts a given ciphertext using a provided private key.
func PEMDecrypt(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	if privateKey == nil {
		return nil, errors.New("private key cannot be nil")
	}
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt ciphertext: %w", err)
	}
	return decrypted, nil
}

// GenerateToken generates a new JWT token with the given payload and secret.
func GenerateToken(payload map[string]interface{}, secret string) (string, error) {
	token, err := GenerateSignedToken(payload, secret)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return token, nil
}

// GenerateSignedToken generates a new JWT token with the given payload and secret.
func GenerateSignedToken(payload map[string]interface{}, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
		"iat":    time.Now().Unix(),
		"sub":    payload["sub"],
		"scopes": payload["scopes"],
	})
	return token.SignedString([]byte(secret))
}

// AuthenticateUser authenticates a user using the provided username and password.
func AuthenticateUser(username string, password string) (string, error) {
	user, err := GetUser(username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}
	return user.JSON()
}

// GetUser retrieves a user by their username.
func GetUser(username string) (*User, error) {
	// Assume we have a database connection and a User struct
	// This is a placeholder, you should replace it with your actual database code
	return &User{
		ID:       "user123",
		Username: username,
		Password: "hashed_password",
	}, nil
}

// User is a struct representing a user.
type User struct {
	ID       string
	Username string
	Password string
}

// JSON returns the user as a JSON string.
func (u *User) JSON() string {
	return fmt.Sprintf(`{"id": "%s", "username": "%s", "password": "%s"}`, u.ID, u.Username, u.Password)
}

func GetHTTPError(status int, message string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       http.NoBody,
		Header:     make(http.Header),
	}
}
```

Please note that you need to replace the `GetUser` function with your actual database code to fetch a user by their username. The `User` struct is also a placeholder and should be replaced with your actual user struct.