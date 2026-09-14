package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"pos-backend/internal/auth/domain"
	"pos-backend/internal/pkg/httputil"
	"pos-backend/internal/pkg/jwtutil"
)

const StaffClaimsKey = "staff_claims"

// AuthMiddleware memvalidasi header Authorization: Bearer <token>
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			httputil.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Header Authorization diperlukan")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			httputil.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Format token harus: Bearer <token>")
			c.Abort()
			return
		}

		claims, err := jwtutil.ValidateToken(parts[1], jwtSecret)
		if err != nil {
			httputil.Error(c, http.StatusUnauthorized, "INVALID_TOKEN", err.Error())
			c.Abort()
			return
		}

		// Simpan data identitas staf ke context Gin
		c.Set(StaffClaimsKey, claims)
		c.Next()
	}
}

// RequireRole membatasi akses endpoint hanya untuk role tertentu
func RequireRole(allowedRoles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := GetStaffClaims(c)
		if !ok {
			httputil.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Identitas staf tidak ditemukan")
			c.Abort()
			return
		}

		for _, r := range allowedRoles {
			if claims.Role == r || claims.Role == domain.RoleOwner { // Owner selalu punya hak akses
				c.Next()
				return
			}
		}

		httputil.Error(c, http.StatusForbidden, "FORBIDDEN", "Anda tidak memiliki izin untuk tindakan ini")
		c.Abort()
	}
}

// GetStaffClaims membaca data identitas staf dari context Gin
func GetStaffClaims(c *gin.Context) (*jwtutil.StaffClaims, bool) {
	val, exists := c.Get(StaffClaimsKey)
	if !exists {
		return nil, false
	}
	claims, ok := val.(*jwtutil.StaffClaims)
	return claims, ok
}
