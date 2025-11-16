package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"pr-reviewer/internal/domain/pullrequest"
	"pr-reviewer/internal/service"
)

type PRHandler struct {
	svc *service.PRService
}

func NewPRHandler(svc *service.PRService) *PRHandler {
	return &PRHandler{svc: svc}
}

type createPRRequest struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

type createPRResponse struct {
	PR *pullrequest.PullRequest `json:"pr"`
}

func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var in createPRRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, 400, "BAD_REQUEST", "invalid json")

		return
	}

	pr, err := h.svc.CreatePR(r.Context(), service.CreatePRInput{
		ID:       in.ID,
		Name:     in.Name,
		AuthorID: in.AuthorID,
	})

	if err != nil {
		switch {
		case errors.Is(err, pullrequest.ErrPRExists):
			writeError(w, 409, "PR_EXISTS", "PR id already exists")
		case errors.Is(err, pullrequest.ErrNotFound):
			writeError(w, 404, "NOT_FOUND", "resource not found")
		default:
			writeError(w, 500, "INTERNAL", err.Error())
		}

		return
	}

	resp := createPRResponse{PR: pr}
	writeJSON(w, http.StatusCreated, resp)
}

type mergePRRequest struct {
	ID string `json:"pull_request_id"`
}

type mergePRResponse struct {
	PR *pullrequest.PullRequest `json:"pr"`
}

func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	var in mergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, 400, "BAD_REQUEST", "invalid json")

		return
	}

	pr, err := h.svc.MergePR(r.Context(), in.ID)
	if err != nil {
		switch {
		case errors.Is(err, pullrequest.ErrNotFound):
			writeError(w, 404, "NOT_FOUND", "resource not found")
		default:
			writeError(w, 500, "INTERNAL", err.Error())
		}

		return
	}

	writeJSON(w, 200, mergePRResponse{PR: pr})
}

type reassignReviewerRequest struct {
	PRID          string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

type reassignReviewerResponse struct {
	PR         *pullrequest.PullRequest    `json:"pr"`
	ReplacedBy string `json:"replaced_by"`
}

func (h *PRHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var in reassignReviewerRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, 400, "BAD_REQUEST", "invalid json")

		return
	}

	pr, replaced, err := h.svc.ReassignReviewer(r.Context(), in.PRID, in.OldReviewerID)
	if err != nil {
		switch {
		case errors.Is(err, pullrequest.ErrNotFound):
			writeError(w, 404, "NOT_FOUND", "resource not found")

		case errors.Is(err, pullrequest.ErrPRMerged):
			writeError(w, 409, "PR_MERGED", "cannot reassign on merged PR")

		case errors.Is(err, pullrequest.ErrNotAssigned):
			writeError(w, 409, "NOT_ASSIGNED", "reviewer is not assigned to this PR")

		case errors.Is(err, pullrequest.ErrNoCandidate):
			writeError(w, 409, "NO_CANDIDATE", "no active replacement candidate in team")

		default:
			writeError(w, 500, "INTERNAL", err.Error())
		}

		return
	}

	writeJSON(w, 200, reassignReviewerResponse{
		PR:         pr,
		ReplacedBy: replaced,
	})
}
