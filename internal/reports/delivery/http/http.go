package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	constants "github.com/pandusatrianura/kasir_api_service/constant"
	"github.com/pandusatrianura/kasir_api_service/internal/reports/service"
	"github.com/pandusatrianura/kasir_api_service/pkg/datetime"
	"github.com/pandusatrianura/kasir_api_service/pkg/response"
)

type ReportHandler struct {
	service service.IReportService
}

func NewReportHandler(service service.IReportService) *ReportHandler {
	return &ReportHandler{
		service: service,
	}
}

// HealthCheck godoc
// @Summary Get health status of reports API
// @Description Get health status of reports API
// @Tags reports
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]string
// @Router /api/reports/health [get]
func (h *ReportHandler) API(w http.ResponseWriter, r *http.Request) {
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

// Report godoc
// @Summary Get sales report
// @Description Get sales report
// @Tags reports
// @Accept json
// @Produce json
// @Param start_date query string false "Start Date"
// @Param end_date query string false "End Date"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/reports [get]
func (h *ReportHandler) Report(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	layout := "2006-01-02"

	if startDate == "" && endDate == "" {
		timeNow := time.Now().UTC()
		date := timeNow.Format(layout)
		startDate = fmt.Sprintf("%s 00:00:00", date)
		endDate = fmt.Sprintf("%s 23:59:59", date)
	} else {
		startDate = fmt.Sprintf("%s 00:00:00", startDate)
		endDate = fmt.Sprintf("%s 23:59:59", endDate)
	}

	startUTC, err := datetime.ParseUTC(startDate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrorCode, constants.ErrReportRequest, err)
		return
	}

	endUTC, err := datetime.ParseUTC(endDate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrorCode, constants.ErrReportRequest, err)
		return
	}

	if startUTC > endUTC {
		response.Error(w, http.StatusBadRequest, constants.ErrorCode, constants.ErrReportRequest, errors.New(constants.ErrStarDate))
		return
	}

	report, err := h.service.Report(startUTC, endUTC)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrorCode, "Report received failed", err)
		return
	}

	response.Success(w, http.StatusCreated, constants.SuccessCode, "Report received successfully", report)
}
