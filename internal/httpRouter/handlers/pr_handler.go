package handlers

import (
	"net/http"
	"pr-reviewer/internal/service"
)

type PRHandler struct {
	svc *service.PRService
}

func NewPRHandler(svc *service.PRService) *PRHandler {
	return &PRHandler{svc: svc}
}

func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request)         {}
func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request)          {}
func (h *PRHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {}
