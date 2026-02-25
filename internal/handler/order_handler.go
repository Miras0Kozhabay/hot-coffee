package handler

import (
	"encoding/json"
	"hot-coffee/internal/logger"
	"hot-coffee/internal/service"
	"hot-coffee/internal/utils"
	"hot-coffee/models"
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
		logger.Log.WithError(err).Error("Order creation failed")
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	logger.Log.WithField("orderID", order.ID).Info("Order created")
	utils.SendJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAll()
	if err != nil {
		logger.Log.WithError(err).Error("Failed to retrieve orders")
		utils.SendError(w, http.StatusInternalServerError, "Failed to retrieve orders")
		return
	}
	utils.SendJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	order, err := h.service.GetByID(id)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to retrieve order")
		utils.SendError(w, http.StatusNotFound, "Order not found")
		return
	}
	logger.Log.WithField("orderID", id).Info("Order retrieved")
	utils.SendJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		// logger.Log.WithError(err).Error("Failed to decode order update payload")
		utils.SendError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	if err := h.service.UpdateOrder(id, order); err != nil {
		logger.Log.WithError(err).Error("Failed to update order")
		if err.Error() == "order not found" {
			// logger.Log.WithField("orderID", id).Error("Order not found for update")
			utils.SendError(w, http.StatusNotFound, err.Error())
		} else {
			// logger.Log.WithField("orderID", id).Error("Order update failed due to bad request")
			utils.SendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	logger.Log.WithField("orderID", id).Info("Order updated")
	utils.SendJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.service.DeleteOrder(id); err != nil {
		logger.Log.WithError(err).Error("Failed to delete order")
		utils.SendError(w, http.StatusNotFound, err.Error())
		return
	}
	logger.Log.WithField("orderID", id).Info("Order deleted")
	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) CloseOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.service.CloseOrder(id); err != nil {
		logger.Log.WithError(err).Error("Failed to close order")
		if err.Error() == "order not found" {
			utils.SendError(w, http.StatusNotFound, err.Error())
		} else {
			logger.Log.WithField("orderID", id).Error("Order closure failed due to bad request")
			utils.SendError(w, http.StatusBadRequest, err.Error())
		}
		return
	}
	logger.Log.WithField("orderID", id).Info("Order closed")
	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Order succesfully closed"})
}
