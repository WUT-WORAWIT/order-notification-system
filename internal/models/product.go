package models

import (
	"time"

	"gorm.io/gorm"
)

// Product defines the structure for product data based on your table schema.
type Product struct {
	ProductID   string    `gorm:"type:varchar(10);primaryKey" json:"product_id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Category    *string   `gorm:"type:varchar(50)" json:"category,omitempty"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	ImageURL    *string   `gorm:"type:varchar(255)" json:"image_url,omitempty"`
	Status      string    `gorm:"type:varchar(10);default:'active'" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for the Product model.
func (Product) TableName() string {
	return "product" // As per your SQL table definition
}

// GetProductByID retrieves a product by its ID from the database.
func GetProductByID(db *gorm.DB, product *Product, productID string) error {
	return db.Where("product_id = ?", productID).First(product).Error
}

// CreateProduct inserts a new product into the database.
func CreateProduct(db *gorm.DB, product *Product) error {
	return db.Create(product).Error
}

// GetAllProducts retrieves all products from the database.
func GetAllProducts(db *gorm.DB, products *[]Product) error {
	return db.Find(products).Error
}
