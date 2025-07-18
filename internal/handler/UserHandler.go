package handler

import (
	entity "PocGo/internal/domain/entities"
	handlerBase "PocGo/internal/handler/base"
	applicationService "PocGo/internal/services"
	setJson "encoding/json"
	httpclient "net/http"
)

type UserHandler interface {
	Update(writer httpclient.ResponseWriter, request *httpclient.Request)
	GetById(writer httpclient.ResponseWriter, request *httpclient.Request)
	GetAll(writer httpclient.ResponseWriter, request *httpclient.Request)
}

type userHandler struct {
	service applicationService.UserService
}

func NewUserHandler(service applicationService.UserService) UserHandler {
	return &userHandler{service: service}
}

func (handler *userHandler) Update(responseWriter httpclient.ResponseWriter, request *httpclient.Request) {
	if err := handlerBase.ValidateHTTPMethod(responseWriter, request, httpclient.MethodPut); err != nil {
		return
	}

	var user entity.User
	if err := setJson.NewDecoder(request.Body).Decode(&user); err != nil {
		handlerBase.SendErrorResponse(responseWriter, err, httpclient.StatusBadRequest)
		return
	}

	if err := handler.service.Update(&user); err != nil {
		handlerBase.SendErrorResponse(responseWriter, err, httpclient.StatusInternalServerError)
		return
	}

	if err := handlerBase.SendJsonResponse(responseWriter, user); err != nil {
		handlerBase.SendErrorResponse(responseWriter, err, httpclient.StatusInternalServerError)
	}
}

func (handler *userHandler) GetById(responseWriter httpclient.ResponseWriter, request *httpclient.Request) {
	if err := handlerBase.ValidationGetMethod(responseWriter, request); err != nil {
		return
	}

	id := handlerBase.GetFromQuery(request, "id")
	if id == "" {
		handlerBase.SendErrorResponse(
			responseWriter,
			setJson.NewDecoder(nil).Decode(nil),
			httpclient.StatusBadRequest,
		)
		return
	}

	user, err := handler.service.GetById(id)
	if err != nil {
		handlerBase.SendErrorResponse(responseWriter, err, httpclient.StatusInternalServerError)
		return
	}

	if err := handlerBase.SendJsonResponse(responseWriter, user); err != nil {
		handlerBase.SendErrorResponse(responseWriter, err, httpclient.StatusInternalServerError)
	}
}

func (handler *userHandler) GetAll(responseWriter httpclient.ResponseWriter, request *httpclient.Request) {
	if err := handlerBase.ValidationGetMethod(responseWriter, request); err != nil {
		return
	}

	date := handlerBase.GetFromQuery(request, "date")

	users, err := handler.service.GetAll(date)
	if err != nil {
		handlerBase.SendErrorResponse(responseWriter, err, httpclient.StatusInternalServerError)
		return
	}

	if err := handlerBase.SendJsonResponse(responseWriter, users); err != nil {
		handlerBase.SendErrorResponse(responseWriter, err, httpclient.StatusInternalServerError)
	}
}
