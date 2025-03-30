package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Id            int     `gorm:"primaryKey;autoIncrement"`
	Name          string  `gorm:"type:varchar(100);not null"`
	Description   string  `gorm:"type:varchar(255);"`
	Price         float64 `gorm:"type:decimal(8,2);not null"`
	StockQuantity int     `gorm:"type:int;not null"`
}
