package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	dispatcher Dispatcher
}

type Dispatcher interface {
	Run(ctx context.Context)
}

func New(dispatcher Dispatcher) *App {
	return &App{
		dispatcher: dispatcher,
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go a.dispatcher.Run(ctx)

	log.Println("worker started")

	<-stop
	log.Println("shutting down...")

	cancel()

	log.Println("shutdown complete")
	return nil
}
