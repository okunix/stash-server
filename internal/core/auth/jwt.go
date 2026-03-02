package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	hmacSecret   = newHmacSecret(32)
	hmacSecretMu = sync.Mutex{}
)

func newHmacSecret(length int) []byte {
	hmacSecretMu.Lock()
	defer hmacSecretMu.Unlock()
	secret := make([]byte, length)
	rand.Read(secret)
	str := hex.EncodeToString(secret)
	fmt.Printf("str: %v\n", str)
	return []byte(str)
}

type UserClaims struct {
	User
	jwt.RegisteredClaims
}

type userClaimsOption func(uc *UserClaims) error

func WithExpirationTime(expiresAt time.Time) userClaimsOption {
	return func(uc *UserClaims) error {
		uc.ExpiresAt = jwt.NewNumericDate(expiresAt)
		return nil
	}
}

func newUserClaims(userID uuid.UUID, opts ...userClaimsOption) (UserClaims, error) {
	userClaims := UserClaims{User: User{UserID: userID}}
	for _, opt := range opts {
		if err := opt(&userClaims); err != nil {
			return userClaims, err
		}
	}
	return userClaims, nil
}

func JWT(userID uuid.UUID, opts ...userClaimsOption) (string, error) {
	claims, err := newUserClaims(userID, opts...)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	hmacSecretMu.Lock()
	defer hmacSecretMu.Unlock()

	return token.SignedString(hmacSecret)
}

func ParseJWT(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (any, error) {
		hmacSecretMu.Lock()
		defer hmacSecretMu.Unlock()
		return hmacSecret, nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok {
		return claims, nil
	}
	return nil, errors.New("failed to parse jwt claims")
}
