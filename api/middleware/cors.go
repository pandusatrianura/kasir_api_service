package middleware

import (
	"net/http"
	"strconv"

	constants "github.com/pandusatrianura/kasir_api_service/constant"
	"github.com/pandusatrianura/kasir_api_service/pkg/response"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var result response.APIResponse

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Method", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "X-API-Key, Content-Type")

		if r.Method == "OPTIONS" {
			result.Code = strconv.Itoa(constants.SuccessCode)
			result.Message = "HTTP Method Options"
			response.WriteJSONResponse(w, http.StatusOK, result)
			return
		}
		next.ServeHTTP(w, r)
	})
}
