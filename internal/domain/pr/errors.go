package pr

import "errors"

var (
	ErrPRExists    = errors.New("pull request already exists")
	ErrPRNotFound  = errors.New("pull request not found")
	ErrPRMerged    = errors.New("pull request already merged")
	ErrNotAssigned = errors.New("user is not assigned to this PR")
	ErrNoCandidate = errors.New("no active replacement candidate available")
)
