package adapters

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Provider struct {
	SecretKey []byte
	Expire    time.Duration
}

func NewProvider(secretKey []byte, expire time.Duration) *Provider {
	return &Provider{
		SecretKey: secretKey,
		Expire:    expire,
	}
}

func (p *Provider) GenerateToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(p.Expire).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(p.SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (p *Provider) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return p.SecretKey, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["sub"].(string), nil
	}
	return "", nil
}
