package user

import "pr-reviewer/internal/domain/pr"

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type ReviewList struct {
	UserID       string                `json:"user_id"`
	PullRequests []pr.PullRequestShort `json:"pull_requests"`
}
