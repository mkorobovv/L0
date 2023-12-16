package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mkorobovv/L0/config"
	"github.com/mkorobovv/L0/internal/models"
)

type DB struct {
	pool *pgxpool.Pool
}

func Connect(cfg *config.Postgres) (*pgxpool.Pool, error) {

	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	return pgxpool.Connect(context.Background(), conn)
}

func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		pool: pool,
	}
}

func (db *DB) CreateTable() error {
	_, err := db.pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS main (
		order_uid VARCHAR(255) PRIMARY KEY,
            track_number VARCHAR(255),
            entry VARCHAR(255),
            delivery_info JSONB,
            payment_info JSONB,
            items JSONB,
            locale VARCHAR(255),
            internal_signature VARCHAR(255),
  				customer_id VARCHAR(255),
  				delivery_service VARCHAR(255),
  				shardkey VARCHAR(255),
  				sm_id INTEGER,
            date_created VARCHAR(255),
            oof_shard VARCHAR(255)
	)
	`)

	return err
}

func (db *DB) SetOrder(order models.OrderJSON) error {
	_, err := db.pool.Exec(context.Background(), `
	INSERT INTO main (order_uid, track_number, entry, delivery_info, payment_info, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`, order.Order_uid, order.Track_number, order.Entry, order.Delivery, order.Payments, order.Items, order.Locale, order.Internal_siganture, order.Customer_id, order.Delivery_service, order.Shardkey, order.Sm_id, order.Date_created, order.OOF_shard)
	return err
}

func (db *DB) GetAllOrders() (orders []models.OrderJSON, err error) {
	res, err := db.pool.Query(context.Background(), `SELECT * FROM main`)

	if err != nil {
		fmt.Printf("Query error: %v", err)
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var executeOrder models.OrderJSON
		err := res.Scan(
			&executeOrder.Order_uid,
			&executeOrder.Track_number,
			&executeOrder.Entry,
			&executeOrder.Delivery,
			&executeOrder.Payments,
			&executeOrder.Items,
			&executeOrder.Locale,
			&executeOrder.Internal_siganture,
			&executeOrder.Customer_id,
			&executeOrder.Delivery_service,
			&executeOrder.Shardkey,
			&executeOrder.Sm_id,
			&executeOrder.Date_created,
			&executeOrder.OOF_shard,
		)

		if err != nil {
			fmt.Printf("Error executing order: %v\n", err)
		}

		orders = append(orders, executeOrder)
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("error executing order: %v", err)
	}

	fmt.Println("Orders downloaded to cache")
	return orders, nil
}

func (db *DB) GetOrderByUID(uid string) (models.OrderJSON, error) {

	var executeOrder models.OrderJSON
	err := db.pool.QueryRow(context.Background(), `SELECT * FROM main WHERE order_uid=$1`, uid).Scan(
		&executeOrder.Order_uid,
		&executeOrder.Track_number,
		&executeOrder.Entry,
		&executeOrder.Delivery,
		&executeOrder.Payments,
		&executeOrder.Items,
		&executeOrder.Locale,
		&executeOrder.Internal_siganture,
		&executeOrder.Customer_id,
		&executeOrder.Delivery_service,
		&executeOrder.Shardkey,
		&executeOrder.Sm_id,
		&executeOrder.Date_created,
		&executeOrder.OOF_shard,
	)

	if err != nil {
		return executeOrder, err
	}

	return executeOrder, nil
}
