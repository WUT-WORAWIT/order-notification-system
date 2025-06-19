package models

import "gorm.io/gorm"

type Order struct {
	ID       uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	ItemCode string  `gorm:"type:varchar" json:"item_code"` // <- เพิ่มตรงนี้
	Item     string  `gorm:"type:varchar" json:"item"`
	Quantity int     `json:"quantity"`
	Price    float64 `gorm:"type:numeric" json:"price"`
	Image    string  `gorm:"type:varchar" json:"image"`
}

func (Order) TableName() string {
	return "orders1"
}

func CreateOrder(db *gorm.DB, order *Order) error {
	result := db.Create(order)
	return result.Error
}

func UpdateOrderStatus(db *gorm.DB, id string, status string) error {
	return db.Model(&Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// GetOrderByID ดึงข้อมูล Order ตาม ID
func GetOrderByID(db *gorm.DB, id string) (*Order, error) {
	var order Order
	result := db.First(&order, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}
