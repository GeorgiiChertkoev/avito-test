package handlers

import (
	"net/http"
	"pr-reviewer/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request)  {}
func (h *UserHandler) GetReviewPRs(w http.ResponseWriter, r *http.Request) {}
