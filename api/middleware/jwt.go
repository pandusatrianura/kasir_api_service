package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	constants "github.com/pandusatrianura/kasir_api_service/constant"
	"github.com/pandusatrianura/kasir_api_service/pkg/response"
	"github.com/spf13/viper"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Roles  string `json:"roles"`
	jwt.RegisteredClaims
}

type JWTConfig struct {
	SecretKey string
	Issuer    string
	Duration  time.Duration
}

func NewJWTConfig(secretKey string, issuer string, durations string) JWTConfig {
	duration, err := time.ParseDuration(durations)
	if err != nil {
		duration = 1 * time.Hour
	}

	return JWTConfig{
		SecretKey: secretKey,
		Issuer:    issuer,
		Duration:  duration,
	}
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var result response.APIResponse
		jwtConfig := NewJWTConfig(viper.GetString("JWT_SECRET_KEY"), viper.GetString("JWT_ISSUER"), viper.GetString("JWT_DURATION"))

		if isPublicRoute(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			result.Code = strconv.Itoa(constants.ErrorCode)
			result.Message = "Empty Authorization Header"
			response.WriteJSONResponse(w, http.StatusInternalServerError, result)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			result.Code = strconv.Itoa(constants.ErrorCode)
			result.Message = "Invalid Authorization Header"
			response.WriteJSONResponse(w, http.StatusInternalServerError, result)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := validateJWT(tokenString, jwtConfig.SecretKey)
		if err != nil {
			result.Code = strconv.Itoa(constants.ErrorCode)
			result.Message = "Invalid Token"
			response.WriteJSONResponse(w, http.StatusInternalServerError, result)
			return
		}

		r.Header.Set("X-User-ID", strconv.Itoa(int(claims.UserID)))
		r.Header.Set("X-User-Email", claims.Email)
		r.Header.Set("X-User-Roles", claims.Roles)

		next.ServeHTTP(w, r)
	})
}

func validateJWT(tokenString, secretKey string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func GenerateJWT(userID uint, email string, roles string, config JWTConfig) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

func isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/health",
		"/login",
	}

	for _, route := range publicRoutes {
		if strings.Contains(path, route) {
			return true
		}
	}

	return false
}
