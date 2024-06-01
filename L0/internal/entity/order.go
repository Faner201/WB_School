package entity

import "time"

type Order struct {
	OrderUID          string    `db:"order_uid" json:"order_uid" validate:"required"`
	TrackNumber       string    `db:"track_number" json:"track_number" validate:"required"`
	Entry             string    `db:"entry" json:"entry" validate:"required"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required,dive,required"`
	Locale            string    `db:"locale" json:"locale" validate:"required"`
	InternalSignature string    `db:"internal_signature" json:"internal_signature"`
	CustomerID        string    `db:"customer_id" json:"customer_id" validate:"required"`
	DeliveryService   string    `db:"delivery_service" json:"delivery_service" validate:"required"`
	Shardkey          string    `db:"shardkey" json:"shardkey" validate:"required"`
	SmID              int       `db:"sm_id" json:"sm_id" validate:"numeric"`
	DateCreated       time.Time `db:"date_created" json:"date_created" validate:"required"`
	OofShard          string    `db:"oof_shard" json:"oof_shard" validate:"required"`
}

type Delivery struct {
	Name    string `db:"name" json:"name" validate:"required"`
	Phone   string `db:"phone" json:"phone" validate:"required"`
	Zip     string `db:"zip" json:"zip" validate:"required"`
	City    string `db:"city" json:"city" validate:"required"`
	Address string `db:"address" json:"address" validate:"required"`
	Region  string `db:"region" json:"region" validate:"required"`
	Email   string `db:"email" json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction  string `db:"transaction" json:"transaction" validate:"required"`
	RequestID    string `db:"request_id" json:"request_id"`
	Currency     string `db:"currency" json:"currency" validate:"required"`
	Provider     string `db:"provider" json:"provider" validate:"required"`
	Amount       int    `db:"amount" json:"amount" validate:"numeric"`
	PaymentDT    int    `db:"payment_dt" json:"payment_dt" validate:"required"`
	Bank         string `db:"bank" json:"bank" validate:"required"`
	DeliveryCost int    `db:"delivery_cost" json:"delivery_cost" validate:"required"`
	GoodsTotal   int    `db:"goods_total" json:"goods_total" validate:"required"`
	CustomFee    int    `db:"custom_fee" json:"custom_fee" validate:"numeric"`
}

type Item struct {
	ChrtID      int    `db:"chrt_id" json:"chrt_id" validate:"numeric"`
	TrackNumber string `db:"track_number" json:"track_number" validate:"required"`
	Price       int    `db:"price" json:"price" validate:"numeric"`
	Rid         string `db:"rid" json:"rid" validate:"required"`
	Name        string `db:"name" json:"name" validate:"required"`
	Sale        int    `db:"sale" json:"sale" validate:"numeric"`
	Size        string `db:"size" json:"size" validate:"required"`
	TotalPrice  int    `db:"total_price" json:"total_price" validate:"numeric"`
	NmID        int    `db:"nm_id" json:"nm_id" validate:"numeric"`
	Brand       string `db:"brand" json:"brand" validate:"required"`
	Status      int    `db:"status" json:"status" validate:"numeric"`
}
