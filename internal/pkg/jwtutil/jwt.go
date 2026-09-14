package jwtutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"pos-backend/internal/auth/domain"
)

var (
	ErrInvalidToken = errors.New("token otentikasi tidak valid atau sudah kadaluwarsa")
)

// StaffClaims merepresentasikan data identitas yang tersimpan di dalam JWT
type StaffClaims struct {
	UserID    uuid.UUID   `json:"user_id"`
	CompanyID uuid.UUID   `json:"company_id"`
	StoreID   uuid.UUID   `json:"store_id"`
	Name      string      `json:"name"`
	Role      domain.Role `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken mencetak JWT digital baru untuk staf yang berhasil login
func GenerateToken(
	user *domain.User,
	storeID uuid.UUID,
	secretKey string,
	duration time.Duration,
) (string, error) {
	claims := StaffClaims{
		UserID:    user.ID,
		CompanyID: user.CompanyID,
		StoreID:   storeID,
		Name:      user.Name,
		Role:      user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ValidateToken membaca dan memverifikasi keaslian signature JWT
func ValidateToken(tokenString string, secretKey string) (*StaffClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &StaffClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*StaffClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
