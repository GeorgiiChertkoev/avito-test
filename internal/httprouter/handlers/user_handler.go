package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"pr-reviewer/internal/domain/pullrequest"
	"pr-reviewer/internal/domain/user"
	"pr-reviewer/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

type setIsActiveRequest struct {
	UserID string `json:"user_id"`
	Active bool   `json:"active"`
}

type setIsActiveResponse struct {
	User *user.User `json:"user"`
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var in setIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, 400, "BAD_REQUEST", "invalid json")

		return
	}

	u, err := h.svc.SetIsActive(r.Context(), in.UserID, in.Active)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			writeError(w, 404, "NOT_FOUND", "resource not found")
		default:
			writeError(w, 500, "INTERNAL", err.Error())
		}

		return
	}

	writeJSON(w, 200, setIsActiveResponse{
		User: u,
	})
}

type getReviewPRsResponse struct {
	UserID       string                         `json:"user_id"`
	PullRequests []pullrequest.PullRequestShort `json:"pull_requests"`
}

func (h *UserHandler) GetReviewPRs(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, 400, "BAD_REQUEST", "`user_id` is required")

		return
	}

	reviews, err := h.svc.GetReviewPRs(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			writeError(w, 404, "NOT_FOUND", "resource not found")
		default:
			writeError(w, 500, "INTERNAL", err.Error())
		}

		return
	}

	writeJSON(w, 200, getReviewPRsResponse{
		UserID:       reviews.UserID,
		PullRequests: reviews.PullRequests,
	})
}
