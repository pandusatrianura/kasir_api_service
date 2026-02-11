package middleware

import (
	"log"
	"net/http"
	"strconv"

	constants "github.com/pandusatrianura/kasir_api_service/constant"
	"github.com/pandusatrianura/kasir_api_service/pkg/response"
)

func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			var result response.APIResponse

			if err := recover(); err != nil {
				log.Printf("Error occurred: %v", err)

				result.Code = strconv.Itoa(constants.ErrorCode)
				result.Message = "Internal Server Error"
				response.WriteJSONResponse(w, http.StatusInternalServerError, result)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
