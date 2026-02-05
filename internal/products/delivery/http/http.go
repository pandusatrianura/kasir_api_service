package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	constants "github.com/pandusatrianura/kasir_api_service/constant"
	"github.com/pandusatrianura/kasir_api_service/internal/products/entity"
	"github.com/pandusatrianura/kasir_api_service/internal/products/service"
	"github.com/pandusatrianura/kasir_api_service/pkg/response"
)

type ProductHandler struct {
	service service.IProductService
}

func NewProductHandler(service service.IProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// HealthCheck godoc
// @Summary Get health status of products API
// @Description Get health status of products API
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]string
// @Router /api/products/health [get]
func (h *ProductHandler) API(w http.ResponseWriter, r *http.Request) {
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

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body entity.RequestProduct true "Product Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, constants.ErrorCode, constants.ErrInvalidMethod, nil)
		return
	}

	var requestProduct entity.RequestProduct
	if err := response.ParseJSON(r, &requestProduct); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrorCode, constants.ErrInvalidProductRequest, err)
		return
	}

	if err := h.service.CreateProduct(&requestProduct); err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrorCode, "Product created failed", err)
		return
	}

	response.Success(w, http.StatusCreated, constants.SuccessCode, "Product created successfully", nil)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body entity.RequestProduct true "Product Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.Error(w, http.StatusMethodNotAllowed, constants.ErrorCode, constants.ErrInvalidMethod, nil)
		return
	}

	var requestProduct entity.RequestProduct

	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrorCode, constants.ErrInvalidProductID, err)
		return
	}

	if err := response.ParseJSON(r, &requestProduct); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrorCode, constants.ErrInvalidProductRequest, err)
		return
	}

	if err := h.service.UpdateProduct(int64(id), &requestProduct); err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrorCode, "Product updated failed", err)
		return
	}

	response.Success(w, http.StatusOK, constants.SuccessCode, "Product updated successfully", nil)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.Error(w, http.StatusMethodNotAllowed, constants.ErrorCode, constants.ErrInvalidMethod, nil)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrorCode, constants.ErrInvalidProductID, err)
		return
	}

	if err := h.service.DeleteProduct(int64(id)); err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrorCode, "Product delete failed", err)
		return
	}

	response.Success(w, http.StatusOK, constants.SuccessCode, "Product deleted successfully", nil)
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Description Get a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/products/{id} [get]
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrorCode, constants.ErrInvalidProductID, err)
		return
	}

	product, err := h.service.GetProductByID(int64(id))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrorCode, "Product retrieved failed", err)
		return
	}

	response.Success(w, http.StatusOK, constants.SuccessCode, "Product retrieved successfully", product)
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Get all products
// @Tags products
// @Accept json
// @Produce json
// @Param name query string false "Product's name"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/products [get]
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	products, err := h.service.GetAllProducts(name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrorCode, "Products retrieved failed", err)
		return
	}

	response.Success(w, http.StatusOK, constants.SuccessCode, "Products retrieved successfully", products)
}
