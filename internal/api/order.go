package api

import (
	"errors"
	"net/http"
	"order-notification-system/internal/models"
	"order-notification-system/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderAPI struct {
	DB *gorm.DB
}

func NewOrderAPI(db *gorm.DB) *OrderAPI {
	return &OrderAPI{DB: db}
}

func (api *OrderAPI) CreateOrder(c *gin.Context) {
	var order models.Order

	// Decode the incoming order request
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if err := models.CreateOrder(api.DB, &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order", "details": err.Error()})
		return
	}

	utils.NotifyNewOrder(strconv.FormatUint(uint64(order.ID), 10), order.ItemCode, order.Item, order.Quantity)

	// ตอบกลับด้วยข้อมูลที่บันทึกสำเร็จ
	c.JSON(http.StatusCreated, order)
}

func (api *OrderAPI) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	var payload struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload", "details": err.Error()})
		return
	}

	if err := models.UpdateOrderStatus(api.DB, orderID, payload.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

// GetProductRequest defines the expected request body for fetching a product.
type GetProductRequest struct {
	ProductID string `json:"product_id" binding:"required"`
}

// GetProduct handles fetching a product by its ID.
// Since the route is POST /getproduct, it expects product_id in the JSON body.
func (api *OrderAPI) GetProduct(c *gin.Context) {
	var req GetProductRequest
	var product models.Product

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// req.ProductID is already validated by `binding:"required"`

	err := models.GetProductByID(api.DB, &product, req.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found", "product_id": req.ProductID})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProducts handles fetching all products.
func (api *OrderAPI) GetProducts(c *gin.Context) {
	var products []models.Product

	err := models.GetAllProducts(api.DB, &products)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products", "details": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusOK, []models.Product{}) // Return empty array if no products found
		return
	}
	c.JSON(http.StatusOK, products)
}

// CreateProduct handles the creation of a new product.
// The route is POST /editproduct, which is a bit unconventional for creation.
func (api *OrderAPI) CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body for creating product", "details": err.Error()})
		return
	}

	// Basic validation can be expanded here (e.g., checking required fields not covered by GORM `not null`)
	if product.ProductID == "" || product.Name == "" || product.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id, name are required, and price must be positive"})
		return
	}

	if err := models.CreateProduct(api.DB, &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}
