package models

type PopularItem struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Count     int    `json:"total_ordered"`
}
