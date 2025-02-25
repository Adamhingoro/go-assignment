package main

import (
	"assignment/core"
	"assignment/cronjob"
	"assignment/database"
	"assignment/handler"
	"assignment/middleware"
	"assignment/model"
	"assignment/router"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/dig"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Received shutdown signal, shutting down...")
		cancel() // Cancel the context
	}()

	env_err := godotenv.Load(".env")
	if env_err != nil {
		if _, ok := os.LookupEnv("DB_HOST"); !ok {
			log.Fatal("DB_HOST environment variable not found")
		}
	}

	container := dig.New()

	// Provide config
	container.Provide(func() database.Config {
		return database.Config{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
		}
	})

	// Provide config
	container.Provide(func() *core.CoreConfig {
		return &core.CoreConfig{
			JWT_KEY: []byte(os.Getenv("JWT_KEY")),
		}
	})

	// Provide dependencies
	container.Provide(database.NewDatabase)
	container.Provide(handler.NewUserHandler)
	container.Provide(handler.NewCompanyHandler)
	container.Provide(handler.NewAuthHandler)
	container.Provide(router.NewRouter)

	// Run database migrations
	container.Invoke(func(db *database.Database) {
		err := db.DB.AutoMigrate(&model.User{}, &model.Company{})
		if err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		log.Println("Database migrated successfully")
		middleware.SetDB(db.DB)
	})

	err := container.Provide(cronjob.NewCheckAndUpdateUsersCronJob)
	if err != nil {
		log.Fatalf("failed to start application: %v", err)
	}

	err = container.Invoke(func(cronScheduler *cron.Cron) {
		log.Println("Cron job started successfully")
	})

	if err != nil {
		log.Fatalf("failed to start application: %v", err)
	}

	// Set up the server and routes
	err = container.Invoke(func(r *router.Router, database *database.Database) {
		log.Println("Registring routes....")
		r.RegisterRoutes()

		go func() {
			log.Println("Server starting on :8080")
			http.ListenAndServe(":8080", r.GetMux())
		}()

		<-ctx.Done()
		log.Println("Closing the database")
		database.CloseDB()
		log.Println("Shutting down the server...")
	})
	if err != nil {
		log.Fatalf("Failed to invoke function: %v", err)
	}

}
