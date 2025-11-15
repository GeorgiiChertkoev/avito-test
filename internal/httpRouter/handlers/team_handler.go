package handlers

import (
	"fmt"
	"net/http"
	"pr-reviewer/internal/service"
)

type TeamHandler struct {
	svc *service.TeamService
}

func NewTeamHandler(svc *service.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {}
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello World")
}
