package main

import (
	"context"
	"fmt"
	"github.com/romakorinenko/hw-test/hw15_go_sql/configs"
	"github.com/romakorinenko/hw-test/hw15_go_sql/internal/dbpool"
	"github.com/romakorinenko/hw-test/hw15_go_sql/internal/handler"
	"github.com/romakorinenko/hw-test/hw15_go_sql/internal/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := configs.MustLoadConfig()
	host := cfg.Server.Host
	port := cfg.Server.Port

	dbPool, err := dbpool.NewDbPool(context.Background(), cfg.Db)
	if err != nil {
		log.Fatalln("cannot create dbPool", err)
	}

	mux := http.NewServeMux()
	userHandler := handler.NewUserHandler(repository.NewUserRepository(dbPool))
	mux.HandleFunc("/users", userHandler.Handle)
	mux.HandleFunc("/users/all", userHandler.GetAll)

	productHandler := handler.NewProductHandler(repository.NewProductRepository(dbPool))
	mux.HandleFunc("/products", productHandler.Handle)
	mux.HandleFunc("/products/all", productHandler.GetAll)

	orderHandler := handler.NewOrderHandler(repository.NewOrderRepository(dbPool))
	mux.HandleFunc("/orders", orderHandler.Handle)
	mux.HandleFunc("/orders/byuserid", orderHandler.GetByUserID)
	mux.HandleFunc("/orders/byuseremail", orderHandler.GetByUserEmail)
	mux.HandleFunc("/orders/user-statistics", orderHandler.GetStatistics)

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", host, port),
		Handler:           mux,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
	}

	go func() {
		_ = server.ListenAndServe()
	}()

	_ = gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("server forced to shutdown:", err)
		return err
	}
	log.Println("server exiting")
	return nil
}
