// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: transfers.sql

package db

import (
	"context"
)

const createTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount
) VALUES ($1, $2, $3) RETURNING id, from_account_id, to_account_id, amount, created_at
`

type CreateTransferParams struct {
	FromAccountID int32   `json:"from_account_id"`
	ToAccountID   int32   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, createTransfer, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAllTransfers = `-- name: DeleteAllTransfers :exec
DELETE FROM transfers
`

func (q *Queries) DeleteAllTransfers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllTransfers)
	return err
}

const getTransferByID = `-- name: GetTransferByID :one
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers WHERE id = $1
`

func (q *Queries) GetTransferByID(ctx context.Context, id int64) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransferByID, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getTransfersByFromAccountID = `-- name: GetTransfersByFromAccountID :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers WHERE from_account_id = $1
`

func (q *Queries) GetTransfersByFromAccountID(ctx context.Context, fromAccountID int32) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, getTransfersByFromAccountID, fromAccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTransfersByToAccountID = `-- name: GetTransfersByToAccountID :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers WHERE to_account_id = $1
`

func (q *Queries) GetTransfersByToAccountID(ctx context.Context, toAccountID int32) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, getTransfersByToAccountID, toAccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTransfers = `-- name: ListTransfers :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers ORDER BY id 
LIMIT $1 OFFSET $2
`

type ListTransfersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, listTransfers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
