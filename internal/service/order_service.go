package service

import (
	"errors"
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/logger"
	"hot-coffee/models"
	"time"
)

type OrderService interface {
	CreateOrder(order *models.Order) error
	GetAll() ([]models.Order, error)
	GetByID(id string) (models.Order, error)
	UpdateOrder(id string, order models.Order) error
	DeleteOrder(id string) error
	CloseOrder(id string) error
}

type orderService struct {
	orderRepo dal.OrderRepository
	menuRepo  dal.MenuRepository
	invRepo   dal.InventoryRepository
}

func NewOrderService(o dal.OrderRepository, m dal.MenuRepository, i dal.InventoryRepository) OrderService {
	return &orderService{
		orderRepo: o,
		menuRepo:  m,
		invRepo:   i,
	}
}

func (s *orderService) CreateOrder(order *models.Order) error {
	if order.CustomerName == "" || len(order.Items) == 0 {
		logger.Log.WithFields(map[string]interface{}{
			"customerName": order.CustomerName,
			"itemsCount":   len(order.Items),
		}).Error("Invalid order: missing customer name or items")
		return errors.New("invalid order: missing  customer name or items")
	}
	requiredIngredients := make(map[string]float64)
	for _, item := range order.Items {
		menuItem, err := s.menuRepo.GetByID(item.ProductID)
		if err != nil {
			logger.Log.WithFields(map[string]interface{}{
				"productID": item.ProductID,
			}).Error("Product not found in menu")
			return fmt.Errorf("product '%s' not found in menu", item.ProductID)
		}
		for _, ing := range menuItem.Ingredients {
			requiredIngredients[ing.IngredientID] += ing.Quantity * float64(item.Quantity)
		}
	}

	inventory, err := s.invRepo.GetAll()
	if err != nil {
		logger.Log.WithError(err).Error("Failed to load inventory")
		return fmt.Errorf("failed to load inventory: %w", err)
	}

	invMap := make(map[string]*models.InventoryItem)
	for i := range inventory {
		invMap[inventory[i].IngredientID] = &inventory[i]
	}

	for ingID, reqQty := range requiredIngredients {
		invItem, exists := invMap[ingID]
		if !exists || invItem.Quantity < reqQty {
			available := 0.0
			name := ingID
			logger.Log.WithFields(map[string]interface{}{
				"ingredientID": ingID,
			}).Info("Ingredient not found in inventory")
			if exists {
				available = invItem.Quantity
				name = invItem.Name
			}
			logger.Log.WithFields(map[string]interface{}{
				"ingredientID": ingID,
				"requiredQty":  reqQty,
				"availableQty": available,
			}).Error("Insufficient stock for ingredient")
			return fmt.Errorf("insufficient stock for '%s': required %.2f, available %.2f", name, reqQty, available)
		}
	}

	for ingID, reqQty := range requiredIngredients {
		logger.Log.WithFields(map[string]interface{}{
			"ingredientID": ingID,
			"requiredQty":  reqQty,
			"availableQty": invMap[ingID].Quantity,
		}).Info("Deducting ingredient from inventory")
		invMap[ingID].Quantity -= reqQty
	}

	if err := s.invRepo.Save(inventory); err != nil {
		logger.Log.WithError(err).Error("Failed to update inventory")
		return errors.New("failed to update inventory records")
	}

	order.ID = fmt.Sprintf("order-%d", time.Now().UnixMilli())
	order.Status = "open"
	order.CreatedAt = time.Now().Format(time.RFC3339)

	orders, _ := s.orderRepo.GetAll()
	orders = append(orders, *order)

	if err := s.orderRepo.SaveAll(orders); err != nil {
		logger.Log.WithError(err).Error("Failed to save order to database")
		return errors.New("failed to save order to database")
	}
	logger.Log.WithFields(map[string]interface{}{
		"orderID":      order.ID,
		"customerName": order.CustomerName,
	}).Info("Order created")
	return nil
}

func (s *orderService) GetAll() ([]models.Order, error) {
	return s.orderRepo.GetAll()
}

func (s *orderService) GetByID(id string) (models.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *orderService) UpdateOrder(id string, updatedOrder models.Order) error {
	// 1. Сначала загружаем текущую версию заказа из базы
	oldOrder, err := s.orderRepo.GetByID(id)
	if err != nil {
		// logger.Log.WithError(err).Error("Failed to retrieve order for update")
		return errors.New("order not found")
	}
	if updatedOrder.Status == "" {
		updatedOrder.Status = oldOrder.Status
	}
	if oldOrder.Status == "closed" {
		// logger.Log.WithFields(map[string]interface{}{
		// 	"orderID": id,
		// }).Error("Attempt to update a closed order")
		return errors.New("cannot modify a closed order")
	}
	if err := s.returnIngredientsToInventory(oldOrder); err != nil {
		logger.Log.WithError(err).Error("Failed to return ingredients to inventory during order update")
		return fmt.Errorf("inventory sync failed (return): %w", err)
	}
	if err := s.deductIngredientsFromInventory(updatedOrder); err != nil {
		_ = s.deductIngredientsFromInventory(oldOrder)
		logger.Log.WithError(err).Error("Failed to deduct ingredients from inventory during order update")
		return fmt.Errorf("inventory sync failed (deduct): %w", err)
	}
	updatedOrder.ID = oldOrder.ID
	updatedOrder.CreatedAt = oldOrder.CreatedAt
	logger.Log.WithFields(map[string]interface{}{
		"orderID": id,
	}).Info("Order updated")
	return s.orderRepo.Update(id, updatedOrder)
}

func (s *orderService) DeleteOrder(id string) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to retrieve order for deletion")
		return errors.New("order not found")
	}
	if order.Status == "closed" {
		logger.Log.WithFields(map[string]interface{}{
			"orderID": id,
		}).Error("Attempt to delete a closed order")
		return errors.New("cannot delete a finalized financial record")
	}
	if order.Status == "open" {
		if err := s.returnIngredientsToInventory(order); err != nil {
			logger.Log.WithError(err).Error("Failed to return ingredients on deletion")
			// slog.Error("Failed to return ingredients on deletion", "error", err)
		}
	}
	logger.Log.WithFields(map[string]interface{}{
		"orderID": id,
	}).Info("Order deleted")
	return s.orderRepo.Delete(id)
}

func (s *orderService) CloseOrder(id string) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}

	if order.Status == "closed" {
		logger.Log.WithFields(map[string]interface{}{
			"orderID": id,
		}).Error("Attempt to close an already closed order")
		return errors.New("order is already closed")
	}

	order.Status = "closed"
	logger.Log.WithFields(map[string]interface{}{
		"orderID": id,
	}).Info("Order closed")
	return s.orderRepo.Update(id, order)
}

func (s *orderService) returnIngredientsToInventory(order models.Order) error {
	inventory, err := s.invRepo.GetAll()
	if err != nil {
		return err
	}
	invMap := make(map[string]*models.InventoryItem)
	for i := range inventory {
		invMap[inventory[i].IngredientID] = &inventory[i]
	}
	for _, item := range order.Items {
		menuItem, err := s.menuRepo.GetByID(item.ProductID)
		if err != nil {
			continue
		}
		for _, ing := range menuItem.Ingredients {
			if invItem, exists := invMap[ing.IngredientID]; exists {
				invItem.Quantity += ing.Quantity * float64(item.Quantity)
			}
		}
	}
	logger.Log.WithFields(map[string]interface{}{
		"orderID": order.ID,
	}).Info("Inventory updated (return) for order")
	return s.invRepo.Save(inventory)
}

func (s *orderService) deductIngredientsFromInventory(order models.Order) error {
	inventory, err := s.invRepo.GetAll()
	if err != nil {
		return err
	}
	invMap := make(map[string]*models.InventoryItem)
	for i := range inventory {
		invMap[inventory[i].IngredientID] = &inventory[i]
	}
	for _, item := range order.Items {
		menuItem, err := s.menuRepo.GetByID(item.ProductID)
		if err != nil {
			logger.Log.WithFields(map[string]interface{}{
				"productID": item.ProductID,
			}).Error("Product not found in menu")
			return fmt.Errorf("товар %s не найден в меню", item.ProductID)
		}

		for _, ing := range menuItem.Ingredients {
			invItem, exists := invMap[ing.IngredientID]
			requiredQty := ing.Quantity * float64(item.Quantity)

			if !exists || invItem.Quantity < requiredQty {
				logger.Log.WithFields(map[string]interface{}{
					"ingredientID": ing.IngredientID,
					"requiredQty":  requiredQty,
					"availableQty": invItem.Quantity,
				}).Error("Insufficient stock for ingredient")
				return fmt.Errorf("недостаточно ингредиента '%s'. требуется: %.2f, в наличии: %.2f",
					ing.IngredientID, requiredQty, invItem.Quantity)
			}
		}
	}
	for _, item := range order.Items {
		menuItem, _ := s.menuRepo.GetByID(item.ProductID)
		for _, ing := range menuItem.Ingredients {
			invMap[ing.IngredientID].Quantity -= ing.Quantity * float64(item.Quantity)
		}
	}
	logger.Log.WithFields(map[string]interface{}{
		"orderID": order.ID,
	}).Info("Inventory updated for order")
	return s.invRepo.Save(inventory)
}
