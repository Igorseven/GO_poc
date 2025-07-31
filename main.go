package main

import (
	bootstrap "PocGo/internal/bootStrap"
	"context"
	logger "log"
	setIO "os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := bootstrap.NewApplication()

	go func() {
		quit := make(chan setIO.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		app.StopScheduler()

		cancel()
	}()

	if err := app.Run(ctx); err != nil {
		logger.Fatal("Erro ao executar aplicação:", err)
	}
}
