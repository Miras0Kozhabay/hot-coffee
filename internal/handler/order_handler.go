package handler

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/internal/utils"
	"hot-coffee/models"
	"log/slog"
	"net/http"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(s service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	if err := h.service.CreateOrder(&order); err != nil {
		slog.Error("Order creation failed", "error", err)
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.SendJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAll()
	if err != nil {
		slog.Error("Failed to retrieve orders", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Failed to retrieve orders")
		return
	}
	utils.SendJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	order, err := h.service.GetByID(id)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Order not found")
		return
	}
	utils.SendJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	if err := h.service.UpdateOrder(id, order); err != nil {
		slog.Error("Failed to update order", "error", err)
		if err.Error() == "order not found" {
			utils.SendError(w, http.StatusNotFound, err.Error())
		} else {
			utils.SendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	utils.SendJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.service.DeleteOrder(id); err != nil {
		slog.Error("Failed to delete order", "error", err)
		utils.SendError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) CloseOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.service.CloseOrder(id); err != nil {
		slog.Error("Failed to close order", "error", err)
		if err.Error() == "order not found" {
			utils.SendError(w, http.StatusNotFound, err.Error())
		} else {
			utils.SendError(w, http.StatusBadRequest, err.Error())
		}
		return
	}
	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Order succesfully closed"})
}
