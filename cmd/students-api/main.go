package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NarsiBhati-codes/students-api/internal/config"
	"github.com/NarsiBhati-codes/students-api/internal/http/handlers/student"
	"github.com/NarsiBhati-codes/students-api/internal/storage/sqlite"
)

func main() {
	// load config 
	cfg := config.MustLoad()

	// logger 
	// database setup 
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	
	// setup router 
	router := http.NewServeMux() 

	router.HandleFunc("POST /api/students", student.New(storage))

	// setup server

	server := http.Server {
		Addr:  cfg.HTTPServer.Addr,
		Handler: router,
	}


	slog.Info("server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1) 
	
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe() 
		if err != nil {
			log.Fatal("failed to start server")
		}
	} ()

	<- done

	slog.Info("shutting down the server")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to Shutdown server", slog.String("error",err.Error()))
	}

	slog.Info("server Shutdown successfully")
}