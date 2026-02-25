package service

import (
	"fmt"

	"hot-coffee/internal/dal"
	"hot-coffee/internal/logger"
	"hot-coffee/models"

	"github.com/sirupsen/logrus"
)

type InventoryService interface {
	AddInventoryItem(item models.InventoryItem) (*models.InventoryItem, error)
	GetAllInventoryItems() ([]models.InventoryItem, error)
	GetInventoryItemByID(id string) (*models.InventoryItem, error)
	UpdateInventoryItem(id string, item models.InventoryItem) (*models.InventoryItem, error)
	DeleteInventoryItem(id string) error
}

type inventoryService struct {
	inventoryRepo dal.InventoryRepository
}

func NewInventoryService(inventoryRepo dal.InventoryRepository) InventoryService {
	return &inventoryService{inventoryRepo: inventoryRepo}
}

func (s *inventoryService) AddInventoryItem(item models.InventoryItem) (*models.InventoryItem, error) {
	if item.IngredientID == "" {
		logger.Log.WithFields(map[string]interface{}{
			"ingredientID": item.IngredientID,
		}).Error("Attempt to add inventory item with missing ingredient ID")
		return nil, fmt.Errorf("ingredient_id is required")
	}
	if item.Name == "" {
		logger.Log.WithFields(map[string]interface{}{
			"ingredientID": item.IngredientID,
		}).Error("Attempt to add inventory item with missing name")
		return nil, fmt.Errorf("name is required")
	}
	if item.Quantity < 0 {
		logger.Log.WithFields(map[string]interface{}{
			"ingredientID": item.IngredientID,
			"quantity":     item.Quantity,
		}).Error("Attempt to add inventory item with negative quantity")
		return nil, fmt.Errorf("quantity cannot be negative")
	}

	existing, err := s.inventoryRepo.GetByID(item.IngredientID)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to check existing inventory item")
		return nil, fmt.Errorf("checking existing item: %w", err)
	}
	if existing != nil {
		logger.Log.WithFields(map[string]interface{}{
			"ingredientID": item.IngredientID,
		}).Error("Attempt to add an inventory item that already exists")
		return nil, fmt.Errorf("inventory item with ID '%s' already exists", item.IngredientID)
	}

	if err := s.inventoryRepo.Create(item); err != nil {
		logger.Log.WithError(err).Error("Failed to create inventory item")
		return nil, fmt.Errorf("creating inventory item: %w", err)
	}

	// logger,Log,WithFields.Info("Inventory item added", "ingredientID", item.IngredientID, "name", item.Name)
	logger.Log.WithFields(logrus.Fields{
		"ingredientID": item.IngredientID,
		"name":         item.Name,
	}).Info("Inventory item added")
	return &item, nil
}

func (s *inventoryService) GetAllInventoryItems() ([]models.InventoryItem, error) {
	items, err := s.inventoryRepo.GetAll()
	if err != nil {
		logger.Log.WithError(err).Error("Failed to retrieve inventory items")
		return nil, fmt.Errorf("fetching inventory items: %w", err)
	}
	return items, nil
}

func (s *inventoryService) GetInventoryItemByID(id string) (*models.InventoryItem, error) {
	item, err := s.inventoryRepo.GetByID(id)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to retrieve inventory item")
		return nil, fmt.Errorf("fetching inventory item: %w", err)
	}
	return item, nil
}

func (s *inventoryService) UpdateInventoryItem(id string, item models.InventoryItem) (*models.InventoryItem, error) {
	existing, err := s.inventoryRepo.GetByID(id)
	if err != nil {
		// logger.Log.WithError(err).Error("Failed to retrieve inventory item for update")
		return nil, fmt.Errorf("fetching inventory item: %w", err)
	}
	if existing == nil {
		return nil, nil
	}

	if item.Name != "" {
		existing.Name = item.Name
	}
	if item.Quantity >= 0 {
		existing.Quantity = item.Quantity
	}
	if item.Unit != "" {
		existing.Unit = item.Unit
	}
	existing.IngredientID = id

	if err := s.inventoryRepo.Update(*existing); err != nil {
		logger.Log.WithError(err).Error("Failed to update inventory item")
		return nil, fmt.Errorf("updating inventory item: %w", err)
	}

	// slog.Info("Inventory item updated", "ingredientID", id)
	logger.Log.WithFields(logrus.Fields{
		"ingredientID": id,
	}).Info("Inventory item updated")
	return existing, nil
}

func (s *inventoryService) DeleteInventoryItem(id string) error {
	existing, err := s.inventoryRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("fetching inventory item: %w", err)
	}
	if existing == nil {
		return nil
	}
	if err := s.inventoryRepo.Delete(id); err != nil {
		return fmt.Errorf("deleting inventory item: %w", err)
	}
	// slog.Info("Inventory item deleted", "ingredientID", id)
	logger.Log.WithFields(logrus.Fields{
		"ingredientID": id,
	}).Info("Inventory item deleted")
	return nil
}
