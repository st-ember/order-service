package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/st-ember/ecommerceprocessor/internal/db"
	"github.com/st-ember/ecommerceprocessor/internal/enum"
	"github.com/st-ember/ecommerceprocessor/internal/model"
)

func BatchStorePurchase(purchases []model.Purchase) error {
	batch := &pgx.Batch{}
	ctx := context.Background()

	for _, purshase := range purchases {
		batch.Queue(
			"INSERT INTO purchase (id, product, customer, purchased_at) VALUES($1, $2, $3, $4)",
			purshase.Id,
			purshase.Product,
			purshase.Customer,
			purshase.PurchasedAt,
		)
	}

	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	results := conn.SendBatch(ctx, batch)
	if results == nil {
		return fmt.Errorf("failed to send batch")
	}

	return nil
}

func PurchaseStatus(id uuid.UUID) (*enum.PurchaseStatus, error) {
	ctx := context.Background()

	var status enum.PurchaseStatus
	query := "SELECT status FROM purchase WHERE id = $1"

	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, query, id)
	err = row.Scan(&status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func DeletePurchase(id uuid.UUID) error {
	ctx := context.Background()

	query := "DELETE FROM purchase WHERE id = $1"

	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	res, err := conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("purchase not found for id: %v", id)
	}

	return nil
}
