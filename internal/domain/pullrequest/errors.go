package pullrequest

import "errors"

var (
	ErrPRExists    = errors.New("pull request already exists")
	ErrPRMerged    = errors.New("pull request already merged")
	ErrNotAssigned = errors.New("user is not assigned to this PR")
	ErrNoCandidate = errors.New("no active replacement candidate available")
	ErrNotFound    = errors.New("resource not found")
)
