package models

type Product struct {
	Id            int     `gorm:"primaryKey;autoIncrement"`
	Name          string  `gorm:"type:varchar(100);not null"`
	Description   string  `gorm:"type:varchar(255);"`
	Price         float32 `gorm:"type:decimal(8,2);not null"`
	StockQuantity int32   `gorm:"type:int;not null"`
}
