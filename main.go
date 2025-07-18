package main

import (
	bootstrap "PocGo/internal/bootStrap"
	logger "log"
	setIO "os"
	"os/signal"
	"syscall"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	// TODO implementar depois, Criando contexto com cancelamento (CancellationToken)
	//ctx, cancel := setContext.WithCancel(setContext.Background())
	//defer cancel()

	app := bootstrap.NewApplication()

	// Configurando graceful shutdown
	go func() {
		quit := make(chan setIO.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		app.StopScheduler()

		//cancel()
	}()

	if err := app.Run(); err != nil {
		logger.Fatal("Erro ao executar aplicação:", err)
	}
}
