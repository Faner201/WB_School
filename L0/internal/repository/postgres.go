package repository

import (
	"L0/internal/entity"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func (r *Repository) InitDB() error {
	db, err := sqlx.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		r.config.DBHost, r.config.DBPort, r.config.DBUser, r.config.DBName, r.config.DBPassword))
	if err != nil {
		log.Err(err).Msg("error open db")
		return err
	}

	if err := db.Ping(); err != nil {
		log.Err(err).Msg("the server is not responding to requests")
		return err
	}

	r.db = db

	return nil
}

func (r *Repository) Close() {
	r.db.Close()
}

func (r *Repository) insertOrder(order entity.Order) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, "+
		"customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES "+
		"($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey, order.SmID,
		order.DateCreated,
		order.SmID)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, "+
		"payment_dt, bank, delivery_cost, goods_total, custom_fee) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec("INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, "+
			"total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, "+
			"$11, $12)",
			order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) сheckOrderDB(order *entity.Order) bool {
	var data entity.Order
	err := r.db.QueryRowx("SELECT * FROM orders WHERE order_uid = $1", order.OrderUID).StructScan(&data)
	return err == nil
}
