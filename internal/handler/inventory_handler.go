package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"hot-coffee/internal/service"
	"hot-coffee/internal/utils"
	"hot-coffee/models"
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
		slog.Warn("Invalid inventory payload", "error", err)
		utils.SendError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	created, err := h.inventoryService.AddInventoryItem(item)
	if err != nil {
		slog.Error("Failed to add inventory item", "error", err)
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.SendJSON(w, http.StatusCreated, created)
}

func (h *InventoryHandler) GetAllInventory(w http.ResponseWriter, r *http.Request) {
	items, err := h.inventoryService.GetAllInventoryItems()
	if err != nil {
		slog.Error("Failed to get inventory items", "error", err)
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
		slog.Error("Failed to get inventory item", "ingredientID", id, "error", err)
		utils.SendError(w, http.StatusInternalServerError, "failed to retrieve inventory item")
		return
	}
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
		slog.Warn("Invalid inventory payload", "error", err)
		utils.SendError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	updated, err := h.inventoryService.UpdateInventoryItem(id, item)
	if err != nil {
		slog.Error("Failed to update inventory item", "ingredientID", id, "error", err)
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
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
		slog.Error("Failed to delete inventory item", "ingredientID", id, "error", err)
		if err.Error() == "inventory item not found" {
			utils.SendError(w, http.StatusNotFound, "inventory item not found")
			return
		}
		utils.SendError(w, http.StatusInternalServerError, "failed to delete inventory item")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
