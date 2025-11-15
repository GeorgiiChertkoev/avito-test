package domain

type CreatePRInput struct {
	ID     string `json:"pull_request_id"`
	Name   string `json:"pull_request_name"`
	Author string `json:"author_id"`
}

type MergePRInput struct {
	ID string `json:"pull_request_id"`
}

type ReassignInput struct {
	PRID      string `json:"pull_request_id"`
	OldUserID string `json:"old_user_id"` // note: OpenAPI says old_user_id in required[]
}
