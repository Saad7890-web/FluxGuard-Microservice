package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID   string   `json:"user_id"`
	Roles    []string `json:"roles"`
	DeviceID string   `json:"device_id"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret []byte
}

func NewManager(secret string) *Manager {
	return &Manager{secret: []byte(secret)}
}

func (m *Manager) Generate(userID string, roles []string, deviceID string, ttl int64) (string, int64, error) {
	exp := time.Now().Add(time.Duration(ttl) * time.Second)

	claims := CustomClaims{
		UserID:   userID,
		Roles:    roles,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", 0, err
	}

	return signed, exp.Unix(), nil
}