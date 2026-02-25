package handler

import (
	"encoding/json"
	"net/http"

	"hot-coffee/internal/logger"
	"hot-coffee/internal/service"
	"hot-coffee/internal/utils"
	"hot-coffee/models"

	"github.com/sirupsen/logrus"
)

type InventoryHandler struct {
	inventoryService service.InventoryService
}

func NewInventoryHandler(inventoryService service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

func (h *InventoryHandler) AddInventory(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		logger.Log.WithError(err).Error("Invalid inventory payload")
		utils.SendError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	created, err := h.inventoryService.AddInventoryItem(item)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to add inventory item")
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	logger.Log.WithFields(logrus.Fields{
		"ingredientID": created.IngredientID,
	}).Info("Inventory item added")
	utils.SendJSON(w, http.StatusCreated, created)
}

func (h *InventoryHandler) GetAllInventory(w http.ResponseWriter, r *http.Request) {
	items, err := h.inventoryService.GetAllInventoryItems()
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get inventory items")
		utils.SendError(w, http.StatusInternalServerError, "failed to retrieve inventory items")
		return
	}

	utils.SendJSON(w, http.StatusOK, items)
}

func (h *InventoryHandler) GetInventoryByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.SendError(w, http.StatusBadRequest, "ingredient ID is required")
		return
	}

	item, err := h.inventoryService.GetInventoryItemByID(id)
	if err != nil {
		logger.Log.WithError(err).
			WithFields(logrus.Fields{
				"ingredientID": id,
			}).
			Error("Failed to get inventory item by ID")
		utils.SendError(w, http.StatusInternalServerError, "failed to retrieve inventory item")
		return
	}
	logger.Log.WithFields(logrus.Fields{
		"ingredientID": id,
	}).Info("Inventory item retrieved by ID")

	utils.SendJSON(w, http.StatusOK, item)
}

func (h *InventoryHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.SendError(w, http.StatusBadRequest, "ingredient ID is required")
		return
	}

	var item models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		// logger.Log.WithError(err).Warn("Invalid inventory payload")
		utils.SendError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	updated, err := h.inventoryService.UpdateInventoryItem(id, item)
	if err != nil {
		logger.Log.WithError(err).
			WithFields(logrus.Fields{
				"ingredientID": id,
			}).
			Error("Failed to update inventory item")
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	logger.Log.WithFields(logrus.Fields{
		"ingredientID": id,
	}).Info("Inventory item updated")
	utils.SendJSON(w, http.StatusOK, updated)
}

func (h *InventoryHandler) DeleteInventory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.SendError(w, http.StatusBadRequest, "ingredient ID is required")
		return
	}
	err := h.inventoryService.DeleteInventoryItem(id)
	if err != nil {
		logger.Log.WithError(err).
			WithFields(logrus.Fields{
				"ingredientID": id,
			}).
			Error("Failed to delete inventory item")
		if err.Error() == "inventory item not found" {
			utils.SendError(w, http.StatusNotFound, "inventory item not found")
			return
		}
		utils.SendError(w, http.StatusInternalServerError, "failed to delete inventory item")
		return
	}
	logger.Log.WithFields(logrus.Fields{
		"ingredientID": id,
	}).Info("Inventory item deleted")
	w.WriteHeader(http.StatusNoContent)
}
