package api

import (
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
