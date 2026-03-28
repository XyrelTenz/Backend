package auth_service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateToken(userID string, role string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey     []byte
	issuer        string
	tokenDuration time.Duration
}

func NewJWTService(secretKey string) JWTService {
	return &jwtService{
		secretKey:     []byte(secretKey),
		issuer:        "backend",
		tokenDuration: 24 * time.Hour,
	}
}

func (s *jwtService) GenerateToken(userID string, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenDuration).Unix(),
			Issuer:    s.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *jwtService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
