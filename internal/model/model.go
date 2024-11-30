package model

import (
	"log"
	"time"
)

type Order struct {
	Order_uid         string    `json:"order_uid" gorm:"primaryKey;unique"`
	Track_number      string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery" gorm:"foreignKey:OrderUID;references:Order_uid"`
	Payment           Payment   `json:"payment" gorm:"foreignKey:OrderUID;references:Order_uid"`
	Items             []Items   `json:"items" gorm:"foreignKey:OrderUID;references:Order_uid"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	Id       int    `json:"-" gorm:"primaryKey"`
	OrderUID string `json:"-" gorm:"index"` // Связь с Order
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Zip      string `json:"zip"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Region   string `json:"region"`
	Email    string `json:"email"`
}

type Payment struct {
	Id           int    `json:"-" gorm:"primaryKey"`
	OrderUID     string `json:"-" gorm:"index"` // Связь с Order
	Transaction  string `json:"transaction" gorm:"primaryKey"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Items struct {
	Id          int    `json:"-" gorm:"primaryKey"`
	OrderUID    string `json:"-" gorm:"index"` // Связь с Order
	ChrtId      int    `json:"chrt_id" gorm:"primaryKey"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func ValidateOrder(order Order) bool {
	if order.Order_uid == "" {
		log.Println("Неверный идентификатор заказа")
		return false
	}
	if len(order.Items) == 0 {
		log.Println("В заказе должно быть хотя бы одно наименование")
		return false
	}
	if order.Payment.Amount <= 0 {
		log.Println(" Сумма платежа должна быть больше нуля")
		return false
	}
	return true
}
