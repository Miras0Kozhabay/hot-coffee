package handler

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/internal/utils"
	"hot-coffee/models"
	"log/slog"
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
		utils.SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.service.CreateMenuItem(menuItem); err != nil {
		slog.Error("Failed to create menu item", "error", err)
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendJSON(w, http.StatusCreated, menuItem)
}

func (h *MenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAll()
	if err != nil {
		slog.Error("Failed to get menu", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	utils.SendJSON(w, http.StatusOK, items)
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.service.GetByID(id)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Menu item not found")
		return
	}
	utils.SendJSON(w, http.StatusOK, item)
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.service.UpdateMenuItem(id, item); err != nil {
		slog.Error("Failed to update menu item", "error", err)
		if err.Error() == "menu item not found" {
			utils.SendError(w, http.StatusNotFound, err.Error())
		} else {
			utils.SendError(w, http.StatusBadRequest, err.Error())
		}
		return
	}
	utils.SendJSON(w, http.StatusOK, item)
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.service.DeleteMenuItem(id); err != nil {
		slog.Error("Failed to delete menu item", "error", err)
		utils.SendError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
