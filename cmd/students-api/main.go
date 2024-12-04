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
)

func main() {
	// load config 
	cfg := config.MustLoad()

	// logger 
	// database setup 
	
	// setup router 
	router := http.NewServeMux() 

	router.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("welcome to students api"))
	})

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

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to Shutdown server", slog.String("error",err.Error()))
	}

	slog.Info("server Shutdown successfully")
}