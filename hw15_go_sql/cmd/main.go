package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romakorinenko/hw-test/hw15_go_sql/internal/config"
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
	cfg := config.MustLoadConfig()
	host := cfg.Server.Host
	port := cfg.Server.Port

	dbPool, err := NewDbPool(context.Background(), cfg.Db)
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

	//orderHandler := handler.NewOrderHandler(repository.NewOrderRepository(dbPool))
	////mux.HandleFunc("/orders", orderHandler.CreateOrderHandler)
	////mux.HandleFunc("/orders", orderHandler.DeleteOrderHandler)
	//mux.HandleFunc("/orders/byuser/{userId}", orderHandler.GetOrdersByUserId)
	//
	//userStatHandler := handler.NewUserStatHandler(repository.NewUserStatRepository(dbPool))
	//mux.HandleFunc("/user-stat/{id}", userStatHandler.GetUserStatById)

	//mux.HandleFunc("POST /order-products/add", nil) // не надо реализовывать, прописать в логике репозитория
	//mux.HandleFunc("DELETE /order-products/remove", nil)

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

func NewDbPool(ctx context.Context, dbCfg *config.Db) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(dbCfg.ConnectionString)
	if err != nil {
		return nil, err
	}
	dbPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}
	return dbPool, nil
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
