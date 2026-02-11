package middleware

import (
	"net/http"
	"strconv"

	constants "github.com/pandusatrianura/kasir_api_service/constant"
	"github.com/pandusatrianura/kasir_api_service/pkg/response"
	"github.com/spf13/viper"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var result response.APIResponse

		// Get API key dari header
		apiKey := r.Header.Get("X-API-Key")

		// Validate
		if apiKey == "" {
			result.Code = strconv.Itoa(constants.ErrorCode)
			result.Message = "API key required"
			response.WriteJSONResponse(w, http.StatusUnauthorized, result)
			return
		}

		if apiKey != viper.GetString("api_key") {
			result.Code = strconv.Itoa(constants.ErrorCode)
			result.Message = "Invalid API key"
			response.WriteJSONResponse(w, http.StatusUnauthorized, result)
			return
		}

		// API key valid, continue
		next.ServeHTTP(w, r)
	})
}
