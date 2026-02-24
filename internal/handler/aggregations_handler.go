package handler

import (
	"log/slog"
	"net/http"

	"hot-coffee/internal/service"
	"hot-coffee/internal/utils"
)

type AggregationHandler struct {
	service service.AggregationService
}

func NewAggregationHandler(s service.AggregationService) *AggregationHandler {
	return &AggregationHandler{service: s}
}

func (h *AggregationHandler) GetTotalSales(w http.ResponseWriter, r *http.Request) {
	total, err := h.service.GetTotalSales()
	if err != nil {
		slog.Error("Failed to calculate total sales", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	response := map[string]float64{
		"total_sales": total,
	}
	utils.SendJSON(w, http.StatusOK, response)
}

func (h *AggregationHandler) GetPopularItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetPopularItems()
	if err != nil {
		slog.Error("Failed to calculate popular items", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	utils.SendJSON(w, http.StatusOK, items)
}
