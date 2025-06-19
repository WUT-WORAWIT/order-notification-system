package handlers

import (
	"errors" // Import errors package
	"fmt"
	"log"      // Import for logging
	"net/http" // Import errors package
	"order-notification-system/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm" // Import gorm for gorm.ErrRecordNotFound
)

// UserHandler holds dependencies for user-related handlers.
type UserHandler struct {
	DB *gorm.DB
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// GetRemark provides a simple remark.
func (h *UserHandler) GetRemark(c *gin.Context) {
	// Base URL, assuming the server runs on localhost:8080
	// This could be made dynamic if needed, e.g., from config or request host.
	baseURL := "http://localhost:8080"

	str := `
Available API Endpoints:

Public Routes:
  General Remark (This Page):
    GET %s/api/
  User Registration:
    POST %s/api/users
      Body (JSON): {"username": "newuser", "password": "password123", "prefix": "Mr.", "first_name": "John", "last_name": "Doe", "email": "john.doe@example.com", "phone_number": "1234567890", "date_of_birth": "YYYY-MM-DD"}
  User Login:
    POST %s/api/login
      Body (JSON): {"username": "existinguser", "password": "password123"}
  Create Order:
    POST %s/order
      Body (JSON): {"item_code": "IC001", "item": "Sample Item", "quantity": 2, "price": 25.50, "image": "http://example.com/image.jpg"}

Protected Routes (Require JWT Bearer Token in 'Authorization' Header, or 'token' query param for WebSocket):
  Get User Profile:
    GET %s/api/profile
  Get User by Username:
    GET %s/api/users/:username (e.g., /api/users/testuser)
  Update User:
    PUT %s/api/users/:username (e.g., /api/users/testuser)
      Body (JSON): {"prefix": "Ms.", "first_name": "Jane"} (fields to update)
  Delete User:
    DELETE %s/api/users/:username (e.g., /api/users/testuser)
  Update Order Status:
    PATCH %s/orders/:id/status (e.g., /orders/1/status)
      Body (JSON): {"status": "Shipped"}
  WebSocket Notifications:
    GET %s/ws?token=YOUR_JWT_TOKEN (Upgrade to WebSocket)
`
	// Format the string with the baseURL
	formattedStr := fmt.Sprintf(str,
		baseURL, // General Remark
		baseURL, // User Registration
		baseURL, // User Login
		baseURL, // Create Order
		baseURL, // Get User Profile
		baseURL, // Get User by Username
		baseURL, // Update User
		baseURL, // Delete User
		baseURL, // Update Order Status
		baseURL, // WebSocket
	)
	c.String(http.StatusOK, formattedStr)
}

// CreateUser handles the creation of a new user.
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON request"})
		return
	}
	fmt.Println(user.Password)
	// Generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate hashed password"})
		return
	}
	user.Password = string(hashedPassword)

	// Create user in database
	err = models.CreateUser(h.DB, &user)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

type TokenStatus struct {
	Code    int
	Message string
}

// GetUserByID handles fetching a user by their username.
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// การเรียก JWTMiddleware ควรทำผ่านการกำหนด route ใน main.go
	// ตัวอย่าง: r.GET("/users/:username", middleware.JWTMiddleware(), handlers.GetUserByID)
	// ที่นี่เราจะสมมติว่า middleware ได้ทำงานแล้วถ้า route ถูกตั้งค่าอย่างถูกต้อง

	// เราจะใช้ c.Param("username") ถ้า route เป็น /users/:username
	// หรือ c.Query("username") ถ้าเป็น /users?username=...
	// จากโค้ดเดิมใช้ c.Query("Username")
	username := c.Query("Username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Username query parameter is required"})
		return
	}

	var user models.User

	err := models.GetUserByID(h.DB, &user, username)
	if err != nil {
		// Use errors.Is for wrapped errors, checking for gorm.ErrRecordNotFound
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
		} else {
			// Log the detailed error on the server for unexpected errors
			log.Printf("Error retrieving user '%s': %v", username, err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve user: " + err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser handles updating user information.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user models.User
	// ควรใช้ path parameter สำหรับระบุ user ที่จะอัปเดต เช่น /users/:username
	// หรือถ้าจะใช้ query parameter ก็ควรเป็น username
	usernameParam := c.Param("username") // สมมติว่า route เป็น /users/:username
	if usernameParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Username path parameter is required"})
		return
	}

	// ดึงข้อมูล user เดิมเพื่อตรวจสอบว่ามีอยู่จริง (optional แต่เป็น good practice)
	var existingUser models.User
	err := models.GetUserByID(h.DB, &existingUser, usernameParam)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User to update not found"})
		return
	}

	// Bind JSON จาก request body เข้าไปยัง user struct ที่จะใช้อัปเดต
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request for update"})
		return
	}

	// ตรวจสอบว่า username ใน payload (ถ้ามี) ตรงกับ username ใน path parameter
	// และตั้งค่า username ให้ถูกต้องก่อน save เพื่อป้องกันการเปลี่ยน username ผ่าน payload โดยไม่ตั้งใจ
	user.Username = usernameParam // Ensure the username from path is used for update

	err = models.UpdateUser(h.DB, &user) // เรียกใช้ models.UpdateUser ที่แก้ไขแล้ว
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update user: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User updated successfully", "data": user})
}

// DeleteUser handles deleting a user.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// ควรใช้ path parameter เช่น /users/:username
	username := c.Param("username") // สมมติว่า route เป็น /users/:username
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Username path parameter is required"})
		return
	}

	// (Optional) ตรวจสอบว่า user มีอยู่จริงหรือไม่ก่อนลบ
	var existingUser models.User
	if err := models.GetUserByID(h.DB, &existingUser, username); errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User to delete not found"})
		return
	}

	err := models.DeleteUser(h.DB, username) // เรียกใช้ models.DeleteUser ที่แก้ไขแล้ว
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete user: " + err.Error()})
		return
	}
	// หากข้อมูลถูกลบสำเร็จ
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User " + username + " has been deleted"})
}
