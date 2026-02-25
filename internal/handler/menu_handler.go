package handler

import (
	"encoding/json"
	"hot-coffee/internal/logger"
	"hot-coffee/internal/service"
	"hot-coffee/internal/utils"
	"hot-coffee/models"
	"net/http"
)

type MenuHandler struct {
	service service.MenuService
}

func NewMenuHandler(s service.MenuService) *MenuHandler {
	return &MenuHandler{service: s}
}

func (h *MenuHandler) CreateMenu(w http.ResponseWriter, r *http.Request) {
	var menuItem models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&menuItem); err != nil {
		logger.Log.WithError(err).Error("Failed to decode menu item creation payload")
		utils.SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.service.CreateMenuItem(menuItem); err != nil {
		logger.Log.WithError(err).Error("Failed to create menu item")
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	logger.Log.WithField("menuItemID", menuItem.ID).Info("Menu item created")
	utils.SendJSON(w, http.StatusCreated, menuItem)
}

func (h *MenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAll()
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get menu")
		utils.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	logger.Log.WithField("itemsCount", len(items)).Info("Menu retrieved")
	utils.SendJSON(w, http.StatusOK, items)
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.service.GetByID(id)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get menu item")
		utils.SendError(w, http.StatusNotFound, "Menu item not found")
		return
	}
	logger.Log.WithField("menuItemID", id).Info("Menu item retrieved")
	utils.SendJSON(w, http.StatusOK, item)
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		logger.Log.WithError(err).Error("Failed to decode menu item update payload")
		utils.SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.service.UpdateMenuItem(id, item); err != nil {
		logger.Log.WithError(err).Error("Failed to update menu item")
		if err.Error() == "menu item not found" {
			// logger.Log.WithField("menuItemID", id).Error("Menu item not found for update")
			utils.SendError(w, http.StatusNotFound, err.Error())
		} else {
			// logger.Log.WithField("menuItemID", id).Error("Menu item update failed due to bad request")
			utils.SendError(w, http.StatusBadRequest, err.Error())
		}
		return
	}
	logger.Log.WithField("menuItemID", id).Info("Menu item updated")
	utils.SendJSON(w, http.StatusOK, item)
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.service.DeleteMenuItem(id); err != nil {
		logger.Log.WithError(err).Error("Failed to delete menu item")
		utils.SendError(w, http.StatusNotFound, err.Error())
		return
	}
	logger.Log.WithField("menuItemID", id).Info("Menu item deleted")
	w.WriteHeader(http.StatusNoContent)
}
