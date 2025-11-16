package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"pr-reviewer/internal/domain/team"
	"pr-reviewer/internal/service"
)

type TeamHandler struct {
	svc *service.TeamService
}

func NewTeamHandler(svc *service.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

type addTeamRequest struct {
	TeamName string        `json:"team_name"`
	Members  []team.Member `json:"members"`
}

type addTeamResponse struct {
	Team *team.Team `json:"team"`
}

func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var in addTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, 400, "BAD_REQUEST", "invalid json")

		return
	}

	t := &team.Team{
		TeamName: in.TeamName,
		Members:  in.Members,
	}

	created, err := h.svc.CreateTeam(r.Context(), t)
	if err != nil {
		switch {
		case errors.Is(err, team.ErrTeamExists):
			writeError(w, 400, "TEAM_EXISTS", "team_name already exists")

		default:
			writeError(w, 500, "INTERNAL", err.Error())
		}

		return
	}

	writeJSON(w, 201, addTeamResponse{
		Team: created,
	})
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		writeError(w, 400, "BAD_REQUEST", "`team_name` is required")

		return
	}

	t, err := h.svc.GetTeam(r.Context(), teamName)
	if err != nil {
		switch {
		case errors.Is(err, team.ErrNotFound):
			writeError(w, 404, "NOT_FOUND", "resource not found")

		default:
			writeError(w, 500, "INTERNAL", err.Error())
		}

		return
	}

	writeJSON(w, 200, t)
}
