package dal

import (
	"errors"
	"hot-coffee/internal/utils"
	"hot-coffee/models"
	"path/filepath"
)

type MenuRepository interface {
	GetAll() ([]models.MenuItem, error)
	GetByID(id string) (models.MenuItem, error)
	Update(id string, item models.MenuItem) error
	Delete(id string) error
	SaveAll(items []models.MenuItem) error
}

type menuRepository struct {
	filePath string
}

func NewMenuRepository(dir string) MenuRepository {
	filePath := filepath.Join(dir, "menu_items.json")
	if err := utils.EnsureFileExists(filePath); err != nil {
		panic("failed to initialize menu database: " + err.Error())
	}
	return &menuRepository{filePath}
}

func (r *menuRepository) GetAll() ([]models.MenuItem, error) {
	return utils.ReadJSON[models.MenuItem](r.filePath)
}

func (r *menuRepository) GetByID(id string) (models.MenuItem, error) {
	items, err := r.GetAll()
	if err != nil {
		return models.MenuItem{}, err
	}

	for _, item := range items {
		if item.ID == id {
			return item, nil
		}
	}
	return models.MenuItem{}, errors.New("menu item not found")
}

func (r *menuRepository) Update(id string, updateItem models.MenuItem) error {
	items, err := r.GetAll()
	if err != nil {
		return err
	}
	found := false
	for i, item := range items {
		if item.ID == id {
			items[i] = updateItem
			found = true
			break
		}
	}
	if !found {
		return errors.New("menu item not found")
	}
	return r.SaveAll(items)
}

func (r *menuRepository) Delete(id string) error {
	items, err := r.GetAll()
	if err != nil {
		return err
	}
	var newItems []models.MenuItem
	found := false
	for _, item := range items {
		if item.ID == id {
			found = true
			continue
		}
		newItems = append(newItems, item)
	}
	if !found {
		return errors.New("menu item not found")
	}
	return r.SaveAll(newItems)
}

func (r menuRepository) SaveAll(items []models.MenuItem) error {
	return utils.WriteJSON[models.MenuItem](r.filePath, items)
}
