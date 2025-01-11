package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Cfg struct {
	Env string
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Sync()

	if err := run(context.TODO(), logger); err != nil {
		logger.Error("Startup", zap.Error(err))
		logger.Sync()
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *zap.Logger) error {
	log := logger.Sugar()

	log.Infow("startup", "main", "started", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	defer log.Infow("shutdown", "main", "completed")

	//----------
	// env config (export APP_ENV=dev)

	var cfg Cfg
	err := envconfig.Process("APP", &cfg)
	if err != nil {
		return err
	}

	//----------
	// http server

	serverErrors := make(chan error, 1)

	// http.ListenAndServe(":8080", http.HandlerFunc(sayHello))
	api := http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(sayHello),
	}

	go func() {
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	// <-shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		// graceful shutdown
		log.Infow("Shutdown", "signal", sig)
		const timeout = 1 * time.Second
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// api's shutdown (graceful)
		err := api.Shutdown(ctx)
		if err != nil {
			api.Close()
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}
		// default:
		// 	time.Sleep(5 * time.Second)
		// 	println("default")
	}

	// print("abc")
	return nil
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	id := rand.Intn(100)
	fmt.Printf("Say hello started %d!\n", id)
	time.Sleep(5 * time.Second)
	fmt.Fprintln(w, "Hello world!", r.Method, r.URL.Path)
	fmt.Println("Say hello finished!", id)
}
