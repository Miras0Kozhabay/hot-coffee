package handler

import "net/http"

func RegisterRoutes(mux *http.ServeMux, invH *InventoryHandler, menuH *MenuHandler, orderH *OrderHandler, agrH *AggregationHandler) {
	//inventory routes
	mux.HandleFunc("POST /inventory", invH.AddInventory)
	mux.HandleFunc("GET /inventory", invH.GetAllInventory)
	mux.HandleFunc("GET /inventory/{id}", invH.GetInventoryByID)
	mux.HandleFunc("PUT /inventory/{id}", invH.UpdateInventory)
	mux.HandleFunc("DELETE /inventory/{id}", invH.DeleteInventory)
	//menu routes
	mux.HandleFunc("POST /menu", menuH.CreateMenu)
	mux.HandleFunc("GET /menu", menuH.GetMenu)
	mux.HandleFunc("GET /menu/{id}", menuH.GetMenuItem)
	mux.HandleFunc("PUT /menu/{id}", menuH.UpdateMenuItem)
	mux.HandleFunc("DELETE /menu/{id}", menuH.DeleteMenuItem)
	// order routes
	mux.HandleFunc("POST /orders", orderH.CreateOrder)
	mux.HandleFunc("GET /orders", orderH.GetOrders)
	mux.HandleFunc("GET /orders/{id}", orderH.GetOrder)
	mux.HandleFunc("PUT /orders/{id}", orderH.UpdateOrder)
	mux.HandleFunc("DELETE /orders/{id}", orderH.DeleteOrder)
	mux.HandleFunc("POST /orders/{id}/close", orderH.CloseOrder)
	//aggregation routes
	mux.HandleFunc("GET /reports/total-sales", agrH.GetTotalSales)
	mux.HandleFunc("GET /reports/popular-items", agrH.GetPopularItems)
}
