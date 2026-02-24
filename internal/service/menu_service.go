package service

import (
	"errors"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type MenuService interface {
	CreateMenuItem(newItem models.MenuItem) error
	GetAll() ([]models.MenuItem, error)
	GetByID(id string) (models.MenuItem, error)
	UpdateMenuItem(id string, item models.MenuItem) error
	DeleteMenuItem(id string) error
}

type menuService struct {
	repo dal.MenuRepository
}

func NewMenuService(repo dal.MenuRepository) MenuService {
	return &menuService{repo: repo}
}

func (s *menuService) CreateMenuItem(newItem models.MenuItem) error {
	if newItem.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if newItem.Name == "" {
		return errors.New("name cannot be empty")
	}
	items, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.ID == newItem.ID {
			return errors.New("menu item with this ID already exists")
		}
	}
	items = append(items, newItem)
	return s.repo.SaveAll(items)
}

func (s *menuService) GetAll() ([]models.MenuItem, error) {
	return s.repo.GetAll()
}

func (s *menuService) GetByID(id string) (models.MenuItem, error) {
	return s.repo.GetByID(id)
}

func (s *menuService) UpdateMenuItem(id string, item models.MenuItem) error {
	if item.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if item.Name == "" {
		return errors.New("name cannot be empty")
	}
	item.ID = id
	return s.repo.Update(id, item)
}

func (s *menuService) DeleteMenuItem(id string) error {
	return s.repo.Delete(id)
}
