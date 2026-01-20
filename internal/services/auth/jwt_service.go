package auth   

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

	"event-management-backend/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	accessSecret string
}

func NewJWTService() *JWTService {
	return &JWTService{
		accessSecret: os.Getenv("JWT_SECRET"),
	}
}

const (
	AccessTTL  = 15 * time.Minute
	RefreshTTL = 7 * 24 * time.Hour
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

func (j *JWTService) GenerateAccessToken(userID uint, role string, permissions []string) (string, error) {
    if j.accessSecret == "" {
        return "", errors.New("jwt access secret missing")
    }

    claims := Claims{
        UserID:      userID,
        Role:        role,
        Permissions: permissions,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(AccessTTL)),
            IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(j.accessSecret))
}

func (j *JWTService) GenerateRefreshToken() (string, string, time.Time, error) {
	b := make([]byte, 40)
	if _, err := rand.Read(b); err != nil {
		return "", "", time.Time{}, err
	}

	raw := base64.URLEncoding.EncodeToString(b)
	hashed := utils.HashToken(raw)
	expiresAt := time.Now().UTC().Add(RefreshTTL)

	return raw, hashed, expiresAt, nil
}

func (j *JWTService) ValidateAccessToken(tokenStr string) (*Claims, error) {
	if tokenStr == "" {
		return nil, errors.New("empty token")
	}

	if j.accessSecret == "" {
		return nil, errors.New("jwt access secret missing")
	}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(j.accessSecret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (j *JWTService) ParseExpiredAccessToken(tokenStr string) (*Claims, error) {
	if tokenStr == "" {
		return nil, errors.New("empty token")
	}
	if j.accessSecret == "" {
		return nil, errors.New("jwt access secret missing")
	}
	parser := jwt.NewParser(
		jwt.WithoutClaimsValidation(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	token, err := parser.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(j.accessSecret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}
