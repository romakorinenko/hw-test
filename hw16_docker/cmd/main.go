package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/romakorinenko/hw-test/hw16_docker/configs"
	"github.com/romakorinenko/hw-test/hw16_docker/internal/dbpool"
	"github.com/romakorinenko/hw-test/hw16_docker/internal/handler"
	"github.com/romakorinenko/hw-test/hw16_docker/internal/repository"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	cfg := configs.MustLoadConfig()
	port := cfg.Server.Port

	dbPool, err := dbpool.NewDBPool(context.Background(), cfg.DB)
	if err != nil {
		log.Fatalln("cannot create dbPool", err)
	}

	db := stdlib.OpenDBFromPool(dbPool)
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatalln("cannot ping database", err)
	}

	goose.SetBaseFS(embedMigrations)
	dialectErr := goose.SetDialect("postgres")
	if dialectErr != nil {
		log.Fatalln("cannot set postgres dialect", dialectErr)
	}
	if migrationsErr := goose.Up(db, "migrations"); migrationsErr != nil {
		log.Fatalln("cannot migrate data", err)
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
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
	}

	log.Println("server started on port: ", port)
	go func() {
		if serverStartErr := server.ListenAndServe(); serverStartErr != nil {
			log.Fatalln("failed to start server", serverStartErr)
		}
	}()

	log.Println(gracefulShutdown(server))
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
