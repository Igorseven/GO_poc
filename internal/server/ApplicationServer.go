package server

import (
	handlers "PocGo/internal/handler"
	"PocGo/internal/middleware"
	applicationService "PocGo/internal/services"
	configIO "fmt"
	muxRouter "github.com/gorilla/mux"
	httpclient "net/http"
)

type ApplicationServer struct {
	userHandler handlers.UserHandler
	router      *muxRouter.Router
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

func (server *ApplicationServer) Start() error {
	configIO.Println("Servidor iniciado na porta 8080")

	handler := middleware.Logging(server.router)

	return httpclient.ListenAndServe(":8080", handler)
}

func (server *ApplicationServer) handleHealth(w httpclient.ResponseWriter, _ *httpclient.Request) {
	w.WriteHeader(httpclient.StatusOK)
}
