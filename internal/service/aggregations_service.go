package service

import (
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type AggregationService interface {
	GetTotalSales() (float64, error)
	GetPopularItems() ([]models.PopularItem, error)
}

type aggregationService struct {
	orderRepo dal.OrderRepository
	menuRepo  dal.MenuRepository
}

func NewAggregationService(orderRepo dal.OrderRepository, menuRepo dal.MenuRepository) AggregationService {
	return &aggregationService{
		orderRepo: orderRepo,
		menuRepo:  menuRepo,
	}
}

func (s *aggregationService) GetTotalSales() (float64, error) {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		return 0, fmt.Errorf("fetching orders: %w", err)

	}

	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return 0, fmt.Errorf("fetching menuu items: %w", err)
	}

	priceMap := make(map[string]float64, len(menuItems))
	for _, m := range menuItems {
		priceMap[m.ID] = m.Price
	}
	var total float64
	for _, order := range orders {
		if order.Status != "closed" {
			continue
		}
		for _, item := range order.Items {
			price, ok := priceMap[item.ProductID]
			if !ok {
				continue
			}
			total += price * float64(item.Quantity)
		}
	}
	return total, nil
}

func (s *aggregationService) GetPopularItems() ([]models.PopularItem, error) {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("fetching orders: %w", err)
	}

	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("fetching menu items: %w", err)
	}

	nameMap := make(map[string]string, len(menuItems))
	for _, m := range menuItems {
		nameMap[m.ID] = m.Name
	}

	countMap := make(map[string]int)
	for _, order := range orders {
		for _, item := range order.Items {
			countMap[item.ProductID] += item.Quantity
		}
	}

	result := make([]models.PopularItem, 0, len(countMap))
	for id, count := range countMap {
		name := nameMap[id]
		if name == "" {
			name = id
		}

		result = append(result, models.PopularItem{
			ProductID: id,
			Name:      name,
			Count:     count,
		})
	}

	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Count > result[i].Count {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result, nil
}
