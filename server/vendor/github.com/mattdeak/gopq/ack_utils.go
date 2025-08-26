package gopq

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	selectItemDetailsQuery  = "SELECT retry_count, ack_deadline FROM %s WHERE id = ?"
	deleteItemQuery         = "DELETE FROM %s WHERE id = ? RETURNING item"
	updateItemForRetryQuery = `
		UPDATE %s 
		SET ack_deadline = ?, retry_count = retry_count + 1
		WHERE id = ?
	`
	expireAckDeadlineQuery = `
		UPDATE %s 
		SET ack_deadline = ?
		WHERE id = ?
	`
)

func nackImpl(ctx context.Context, db *sql.DB, tableName string, id int64, opts AckOpts) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback() // will fail if committed, but that's fine
	}()

	var retryCount int
	var ackDeadline int64
	err = tx.QueryRow(fmt.Sprintf(selectItemDetailsQuery, tableName), id).Scan(&retryCount, &ackDeadline)
	if err != nil {
		return fmt.Errorf("failed to get item details: %w", err)
	}

	// Check if the ack deadline has expired
	if ackDeadline < time.Now().Unix() {
		return fmt.Errorf("ack deadline has expired, cannot nack")
	}

	// Check if we have reached the maximum number of retries
	if retryCount >= opts.MaxRetries && opts.MaxRetries != InfiniteRetries {
		return handleTooManyRetries(tx, tableName, id, opts)
	}

	// Use the maximum of retryBackoff and ackTimeout
	newDeadline := time.Now().Add(max(opts.RetryBackoff, opts.AckTimeout)).Unix()
	_, err = tx.Exec(fmt.Sprintf(updateItemForRetryQuery, tableName), newDeadline, id)
	if err != nil {
		return fmt.Errorf("failed to update item for retry: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func handleTooManyRetries(tx *sql.Tx, tableName string, id int64, opts AckOpts) error {
	var item []byte
	err := tx.QueryRow(fmt.Sprintf(deleteItemQuery, tableName), id).Scan(&item)
	if err != nil {
		return fmt.Errorf("failed to delete item for on failure: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	if len(opts.FailureCallbacks) > 0 {
		for _, fn := range opts.FailureCallbacks {
			err := fn(Msg{
				ID:   id,
				Item: item,
			})
			if err != nil {
				return fmt.Errorf("failed to execute failure callback: %w", err)
			}
		}
	}

	return nil
}

// max returns the maximum of two time.Duration values
func max(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

func expireAckDeadline(db *sql.DB, name string, id int64) error {
	// expiredTime is 1 second in the past to ensure that the ack deadline is expired
	expiredTime := time.Now().Add(-1 * time.Second).Unix()
	_, err := db.Exec(fmt.Sprintf(expireAckDeadlineQuery, name), expiredTime, id)
	return err
}
