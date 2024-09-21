package handlers

import (
	"bookstore_api/internal/services"
	"net/http"
)

type UserHandler struct {
	*Handler
	userService *services.UserService
}

func NewUserHandler(handler *Handler, userService *services.UserService) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

}
