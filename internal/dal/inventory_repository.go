package dal

import (
	"fmt"
	"path/filepath"

	"hot-coffee/internal/logger"
	"hot-coffee/internal/utils"
	"hot-coffee/models"
)

type InventoryRepository interface {
	GetAll() ([]models.InventoryItem, error)
	GetByID(id string) (*models.InventoryItem, error)
	Create(item models.InventoryItem) error
	Update(item models.InventoryItem) error
	Delete(id string) error
	Save(items []models.InventoryItem) error
}

type inventoryRepository struct {
	filePath string
}

func NewInventoryRepository(dir string) InventoryRepository {
	filePath := filepath.Join(dir, "inventory.json")
	if err := utils.EnsureFileExists(filePath); err != nil {
		logger.Log.WithError(err).Fatal("Failed to initialize inventory database")
	}
	return &inventoryRepository{filePath}
}

func (r *inventoryRepository) GetAll() ([]models.InventoryItem, error) {
	items, err := utils.ReadJSON[models.InventoryItem](r.filePath)
	if err != nil {
		return nil, fmt.Errorf("reading inventory: %w", err)
	}
	return items, nil
}

func (r *inventoryRepository) GetByID(id string) (*models.InventoryItem, error) {
	items, err := r.GetAll()
	if err != nil {
		return nil, err
	}
	for i := range items {
		if items[i].IngredientID == id {
			return &items[i], nil
		}
	}
	return nil, nil
}

func (r *inventoryRepository) Create(item models.InventoryItem) error {
	items, err := r.GetAll()
	if err != nil {
		return err
	}
	items = append(items, item)
	logger.Log.WithField("item", item).Info("Creating inventory item")
	return r.Save(items)
}

func (r *inventoryRepository) Update(item models.InventoryItem) error {
	items, err := r.GetAll()
	if err != nil {
		return err
	}
	for i := range items {
		if items[i].IngredientID == item.IngredientID {
			items[i] = item
			return r.Save(items)
		}
	}

	return fmt.Errorf("inventory item not found: %s", item.IngredientID)
}

func (r *inventoryRepository) Delete(id string) error {
	items, err := r.GetAll()
	if err != nil {
		return err
	}
	newItems := make([]models.InventoryItem, 0, len(items))
	found := false
	for _, it := range items {
		if it.IngredientID == id {
			found = true
			continue
		}
		newItems = append(newItems, it)
	}
	if !found {
		return fmt.Errorf("inventory item not found: %s", id)
	}
	return r.Save(newItems)
}

func (r *inventoryRepository) Save(items []models.InventoryItem) error {
	if err := utils.WriteJSON(r.filePath, items); err != nil {
		return fmt.Errorf("writing inventory: %w", err)
	}
	return nil
}
