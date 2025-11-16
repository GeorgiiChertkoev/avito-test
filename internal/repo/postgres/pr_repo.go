package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"pr-reviewer/internal/domain/pullrequest"
)

type PRRepo struct {
	tx *sql.Tx
}

func NewPRRepo(tx *sql.Tx) *PRRepo {
	return &PRRepo{tx: tx}
}

func (r *PRRepo) CreatePR(ctx context.Context, pr *pullrequest.PullRequest) error {
	// Check if exists
	var exists int
	err := r.tx.QueryRowContext(ctx,
		`SELECT 1 FROM pull_request WHERE pull_request_id = $1`,
		pr.ID,
	).Scan(&exists)

	if err == nil {
		return pullrequest.ErrPRExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	query := `
		INSERT INTO pull_request (
			pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = r.tx.ExecContext(ctx, query,
		pr.ID,
		pr.Name,
		pr.AuthorID,
		pr.Status,
		pr.CreatedAt,
		pr.MergedAt,
	)
	if err != nil {
		return err
	}

	// reviewers
	for _, reviewerID := range pr.AssignedReviewers {
		_, err := r.tx.ExecContext(ctx,
			`INSERT INTO pull_request_reviewer (pull_request_id, reviewer_id)
			 VALUES ($1, $2)`,
			pr.ID, reviewerID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PRRepo) GetPR(ctx context.Context, id string) (*pullrequest.PullRequest, error) {
	query := `
		SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		FROM pull_request
		WHERE pull_request_id = $1
	`

	pr := &pullrequest.PullRequest{}
	err := r.tx.QueryRowContext(ctx, query, id).Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, pullrequest.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// reviewers
	rows, err := r.tx.QueryContext(ctx,
		`SELECT reviewer_id FROM pull_request_reviewer WHERE pull_request_id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	var reviewers []string
	for rows.Next() {
		var rid string
		if err := rows.Scan(&rid); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, rid)
	}
	pr.AssignedReviewers = reviewers

	return pr, nil
}

func (r *PRRepo) UpdatePR(ctx context.Context, pr *pullrequest.PullRequest) error {
	// Load existing PR
	existing, err := r.GetPR(ctx, pr.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return pullrequest.ErrNotFound
	}

	// Cannot update merged PR
	if existing.Status == pullrequest.StatusMerged {
		return pullrequest.ErrPRMerged
	}

	query := `
		UPDATE pull_request
		SET pull_request_name = $1,
		    author_id = $2,
		    status = $3,
		    created_at = $4,
		    merged_at = $5
		WHERE pull_request_id = $6
	`

	_, err = r.tx.ExecContext(ctx, query,
		pr.Name,
		pr.AuthorID,
		pr.Status,
		pr.CreatedAt,
		pr.MergedAt,
		pr.ID,
	)
	if err != nil {
		return err
	}

	// reset reviewers
	_, err = r.tx.ExecContext(ctx,
		`DELETE FROM pull_request_reviewer WHERE pull_request_id = $1`,
		pr.ID,
	)
	if err != nil {
		return err
	}

	for _, rid := range pr.AssignedReviewers {
		_, err := r.tx.ExecContext(ctx,
			`INSERT INTO pull_request_reviewer (pull_request_id, reviewer_id)
			 VALUES ($1, $2)`,
			pr.ID, rid,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PRRepo) GetReviewerPRs(ctx context.Context, userID string) ([]pullrequest.PullRequestShort, error) {
	query := `
		SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
		FROM pull_request pr
		JOIN pull_request_reviewer r
		    ON pr.pull_request_id = r.pull_request_id
		WHERE r.reviewer_id = $1
		ORDER BY pr.created_at DESC
	`

	rows, err := r.tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	var result []pullrequest.PullRequestShort

	for rows.Next() {
		var pr pullrequest.PullRequestShort
		if err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		result = append(result, pr)
	}

	return result, nil
}
