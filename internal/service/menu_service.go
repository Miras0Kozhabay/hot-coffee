package service

import (
	"errors"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/logger"
	"hot-coffee/models"

	"github.com/sirupsen/logrus"
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
		logger.Log.WithFields(logrus.Fields{
			"menuItemID": newItem.ID,
			"price":      newItem.Price,
		}).Error("Menu item price must be greater than zero")
		return errors.New("price must be greater than zero")
	}
	if newItem.Name == "" {
		logger.Log.WithFields(logrus.Fields{
			"menuItemID": newItem.ID,
		}).Error("Menu item name cannot be empty")
		return errors.New("name cannot be empty")
	}
	if newItem.ID == "" {
		logger.Log.WithFields(logrus.Fields{
			"name": newItem.Name,
		}).Error("Menu item ID cannot be empty")
		return errors.New("ID cannot be empty")
	}
	items, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.ID == newItem.ID {
			logger.Log.WithFields(logrus.Fields{
				"menuItemID": newItem.ID,
			}).Error("Menu item with this ID already exists")
			return errors.New("menu item with this ID already exists")
		}
	}
	items = append(items, newItem)
	logger.Log.WithFields(logrus.Fields{
		"menuItemID": newItem.ID,
		"name":       newItem.Name,
	}).Info("Menu item created")
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
		logger.Log.WithFields(logrus.Fields{
			"menuItemID": id,
			"price":      item.Price,
		}).Error("Menu item price must be greater than zero")
		return errors.New("price must be greater than zero")
	}
	if item.Name == "" {
		logger.Log.WithFields(logrus.Fields{
			"menuItemID": id,
		}).Error("Menu item name cannot be empty")
		return errors.New("name cannot be empty")
	}
	item.ID = id
	logger.Log.WithFields(logrus.Fields{
		"menuItemID": id,
		"name":       item.Name,
	}).Info("Menu item updated")
	return s.repo.Update(id, item)
}

func (s *menuService) DeleteMenuItem(id string) error {
	logger.Log.WithFields(logrus.Fields{
		"menuItemID": id,
	}).Info("Menu item deleted")
	logger.Log.WithFields(logrus.Fields{
		"menuItemID": id,
	}).Info("Menu item deleted")
	return s.repo.Delete(id)
}
