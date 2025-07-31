package server

import (
	handlers "PocGo/internal/handler"
	"PocGo/internal/middleware"
	applicationService "PocGo/internal/services"
	"context"
	"errors"
	configIO "fmt"
	muxRouter "github.com/gorilla/mux"
	httpclient "net/http"
	"time"
)

type ApplicationServer struct {
	userHandler handlers.UserHandler
	router      *muxRouter.Router
	httpServer  *httpclient.Server
}

func NewServer(
	services *applicationService.Services,
) *ApplicationServer {
	server := &ApplicationServer{
		userHandler: handlers.NewUserHandler(services.User),
		router:      muxRouter.NewRouter(),
	}

	server.setupRoutes()
	return server
}

func (server *ApplicationServer) setupRoutes() {

	server.router = muxRouter.NewRouter()

	server.router.Handle("/user/get_user_by_id",
		middleware.Logging(httpclient.HandlerFunc(server.userHandler.GetById))).
		Methods(httpclient.MethodGet).
		Queries("id", "{id}")

	server.router.Handle("/user/get_all_users",
		middleware.Logging(httpclient.HandlerFunc(server.userHandler.GetAll))).
		Methods(httpclient.MethodGet).
		Queries("data", "{data}")

	server.router.Handle("/user/update_user",
		middleware.Logging(httpclient.HandlerFunc(server.userHandler.Update))).
		Methods(httpclient.MethodPut)

	server.router.Handle("/health",
		middleware.Logging(httpclient.HandlerFunc(server.handleHealth))).
		Methods(httpclient.MethodGet)

}

func (server *ApplicationServer) Start(ctx context.Context) error {
	configIO.Println("Servidor iniciado na porta 8080")

	handler := middleware.Logging(server.router)

	server.httpServer = &httpclient.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil && !errors.Is(err, httpclient.ErrServerClosed) {
			configIO.Printf("Erro ao iniciar servidor: %v\n", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return server.httpServer.Shutdown(shutdownCtx)
}

func (server *ApplicationServer) Shutdown(ctx context.Context) error {
	if server.httpServer != nil {
		return server.httpServer.Shutdown(ctx)
	}
	return nil
}

func (server *ApplicationServer) handleHealth(w httpclient.ResponseWriter, _ *httpclient.Request) {
	w.WriteHeader(httpclient.StatusOK)
}
