package dal

import (
	"errors"
	"hot-coffee/internal/logger"
	"hot-coffee/internal/utils"
	"hot-coffee/models"
	"path/filepath"
)

type OrderRepository interface {
	GetAll() ([]models.Order, error)
	GetByID(id string) (models.Order, error)
	SaveAll(orders []models.Order) error
	Update(id string, order models.Order) error
	Delete(id string) error
}
type orderRepository struct {
	filePath string
}

func NewOrderRepository(dir string) OrderRepository {
	filePath := filepath.Join(dir, "orders.json")
	if err := utils.EnsureFileExists(filePath); err != nil {
		logger.Log.WithError(err).Fatal("Failed to initialize order database")
	}
	return &orderRepository{filePath}
}

func (r *orderRepository) GetAll() ([]models.Order, error) {
	return utils.ReadJSON[models.Order](r.filePath)
}

func (r *orderRepository) GetByID(id string) (models.Order, error) {
	orders, err := r.GetAll()
	if err != nil {
		return models.Order{}, err
	}
	for _, order := range orders {
		if order.ID == id {
			return order, nil
		}
	}
	return models.Order{}, errors.New("order not found")
}

func (r *orderRepository) SaveAll(orders []models.Order) error {
	return utils.WriteJSON[models.Order](r.filePath, orders)
}

func (r *orderRepository) Update(id string, updatedOrder models.Order) error {
	orders, err := r.GetAll()
	if err != nil {
		return err
	}
	found := false
	for i, order := range orders {
		if order.ID == id {
			orders[i] = updatedOrder
			found = true
			break
		}
	}
	if !found {
		return errors.New("order not found")
	}
	return r.SaveAll(orders)
}

func (r *orderRepository) Delete(id string) error {
	orders, err := r.GetAll()
	if err != nil {
		return err
	}
	var newOrders []models.Order
	found := false
	for _, order := range orders {
		if order.ID == id {
			found = true
			continue
		}
		newOrders = append(newOrders, order)
	}
	if !found {
		return errors.New("order not found")
	}
	return r.SaveAll(newOrders)
}
