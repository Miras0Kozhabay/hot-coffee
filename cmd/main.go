package main

import (
	"context"
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
	"hot-coffee/internal/utils"
	"log"
	"net/http"
	"os"
	"log/slog"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := utils.Load()
	 if err != nil {
		slog.Error("Failed to load configurations", "error", err)
	 }

	invRepo := dal.NewInventoryRepository(cfg.DataDir)
	menuRepo := dal.NewMenuRepository(cfg.DataDir)
	orderRepo := dal.NewOrderRepository(cfg.DataDir)

	invService := service.NewInventoryService(invRepo)
	menuService := service.NewMenuService(menuRepo)
	orderService := service.NewOrderService(orderRepo, menuRepo, invRepo)
	agrService := service.NewAggregationService(orderRepo, menuRepo)

	invHandler := handler.NewInventoryHandler(invService)
	menuHandler := handler.NewMenuHandler(menuService)
	orderHandler := handler.NewOrderHandler(orderService)
	agrHandler := handler.NewAggregationHandler(agrService)
	// Настройка роутера 
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, invHandler, menuHandler, orderHandler, agrHandler)
	// Запуск серверa
	fileServer := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fileServer)
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}
	go func() {
		fmt.Printf("starting server on port %s...\n", cfg.Port)
		fmt.Printf("dir is %s...\n", cfg.DataDir)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("\nshutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// fmt.Errorf("server forced to shutdown: %w", err) заменить библиотекой лог
		
	}
	fmt.Println("server exiting")
}
