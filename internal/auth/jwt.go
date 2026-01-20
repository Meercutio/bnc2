package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID      string `json:"uid"`
	DisplayName string `json:"name,omitempty"`
	jwt.RegisteredClaims
}

// Service is a small auth helper that encapsulates a signing secret.
// It makes it easier to inject auth into handlers without passing raw []byte everywhere.
type Service struct {
	secret []byte
}

func NewService(secret []byte) *Service {
	// Copy to avoid accidental mutation from outside.
	s := make([]byte, len(secret))
	copy(s, secret)
	return &Service{secret: s}
}

func (s *Service) Sign(userID string, ttl time.Duration) (string, error) {
	return Sign(s.secret, userID, "", ttl)
}

// SignWithName signs a token and embeds a displayName into claims.
// Useful when you want to render nicknames client-side without an extra DB call.
func (s *Service) SignWithName(userID, displayName string, ttl time.Duration) (string, error) {
	return Sign(s.secret, userID, displayName, ttl)
}

func (s *Service) Verify(token string) (*Claims, error) {
	return Verify(s.secret, token)
}

func Sign(secret []byte, userID, displayName string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID:      userID,
		DisplayName: displayName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

func Verify(secret []byte, token string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
