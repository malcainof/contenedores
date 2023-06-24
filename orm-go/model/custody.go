package model

import "github.com/shopspring/decimal"

type Custody struct {
	Period   string          `gorm:"primary_key;column:period;type:TEXT NOT NULL" json:"period"`
	Stock    string          `gorm:"primary_key;column:stock;type:TEXT NOT NULL" json:"stock"`
	ClientId string          `gorm:"primary_key;column:client_idd;type:TEXT NOT NULL" json:"client_id"`
	Market   string          `gorm:"column:market;type:TEXT NOT NULL" json:"market"`
	Price    decimal.Decimal `gorm:"column:price;type:decimal(20, 6) NOT NULL" json:"price"`
	Quantity int32           `gorm:"column:quantity;type:INT4 NOT NULL" json:"quantity"`
}

func (t *Custody) TableName() string {
	return "custody"
}
