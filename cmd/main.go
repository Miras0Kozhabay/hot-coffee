package main

import (
	"context"
	"flag"
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
	"hot-coffee/internal/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Парсим флаги 
	port := flag.String("port", "8080", "")
	dir := flag.String("dir", "data", "")
	flag.Usage = func() {
		utils.PrintHelp()
	}
	flag.Parse()

	invRepo := dal.NewInventoryRepository(*dir)
	menuRepo := dal.NewMenuRepository(*dir)
	orderRepo := dal.NewOrderRepository(*dir)

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
		Addr:    ":" + *port,
		Handler: mux,
	}
	go func() {
		fmt.Printf("starting server on port %s...\n", *port)
		fmt.Printf("dir is %s...\n", *dir)
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
