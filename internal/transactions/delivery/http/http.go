package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	constants "github.com/pandusatrianura/kasir_api_service/constant"
	"github.com/pandusatrianura/kasir_api_service/internal/transactions/entity"
	"github.com/pandusatrianura/kasir_api_service/internal/transactions/service"
	"github.com/pandusatrianura/kasir_api_service/pkg/response"
)

type TransactionHandler struct {
	service service.ITransactionsService
}

func NewTransactionsHandler(service service.ITransactionsService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// HealthCheck godoc
// @Summary Get health status of transactions/checkout API
// @Description Get health status of transactions/checkout API
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]string
// @Router /api/transactions/health [get]
func (h *TransactionHandler) API(w http.ResponseWriter, r *http.Request) {
	var result response.APIResponse
	svcHealthCheckResult := h.service.API()

	if svcHealthCheckResult.IsHealthy {
		result.Code = strconv.Itoa(constants.SuccessCode)
		result.Message = fmt.Sprintf("%s is healthy", svcHealthCheckResult.Name)
		response.WriteJSONResponse(w, http.StatusOK, result)
		return
	}

	result.Code = strconv.Itoa(constants.ErrorCode)
	result.Message = fmt.Sprintf("%s is not healthy", svcHealthCheckResult.Name)
	response.WriteJSONResponse(w, http.StatusServiceUnavailable, result)
	return
}

// Checkout godoc
// @Summary Checkout products
// @Description Checkout products
// @Tags transactions
// @Accept json
// @Produce json
// @Param X-API-Key header string true "your-secret-api-key-here"
// @Param checkout body entity.Checkout true "Checkout Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/transactions/checkout [post]
func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, constants.ErrorCode, constants.ErrInvalidMethod, nil)
		return
	}

	var (
		request   entity.Checkout
		checkouts []entity.CheckoutRequest
		err       error
		resp      interface{}
	)

	if err = response.ParseJSON(r, &request); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrorCode, constants.ErrInvalidCheckoutRequest, err)
		return
	}

	if len(request.Checkouts) == 0 {
		response.Error(w, http.StatusBadRequest, constants.ErrorCode, constants.ErrInvalidCheckoutRequest, errors.New("checkouts is empty"))
		return
	}

	for _, checkout := range request.Checkouts {
		checkouts = append(checkouts, checkout)
	}

	if resp, err = h.service.Checkout(checkouts); err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrorCode, "Checkout created failed", err)
		return
	}

	if resp == nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrorCode, "Checkout created failed", errors.New("checkout response is nil"))
		return
	}

	response.Success(w, http.StatusCreated, constants.SuccessCode, "Checkout created successfully", resp)
}
