package handlers

import (
	"errors"
	"log"
	"net/http"                                // Import errors package
	"order-notification-system/internal/auth" // Import the new auth package
	"order-notification-system/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm" // Import gorm
)

// ProfileHandler holds dependencies for profile-related handlers.
type ProfileHandler struct {
	DB *gorm.DB
}

// NewProfileHandler creates a new ProfileHandler instance.
func NewProfileHandler(db *gorm.DB) *ProfileHandler {
	return &ProfileHandler{DB: db}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Get claims from middleware
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "No token claims found",
		})
		return
	}

	// Type assert claims
	customClaims, ok := claims.(*auth.CustomClaims)
	if !ok {
		// Log the actual type of claims for debugging if assertion fails
		log.Printf("Error: could not assert claims to *auth.CustomClaims. Actual type: %T. Value: %+v", claims, claims)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to parse token claims to expected type",
		})
		return
	}

	// Get user from database using the username from claims
	var user models.User
	// Use customClaims.Username (or customClaims.Subject which should be the same)
	if err := models.GetUserByID(h.DB, &user, customClaims.Username); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // Use errors.Is for wrapped errors
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
		} else {
			log.Printf("Error retrieving user profile for '%s': %v", customClaims.Username, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to retrieve user profile: " + err.Error(),
			})
		}
		return
	}
	// Return user profile (excluding sensitive data)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"username":   user.Username,
			"prefix":     user.Prefix,
			"first_name": user.First_name,
			"last_name":  user.Last_name,
			"email":      user.Email,
			// "phone_number": user.PhoneNumber, // Consider if phone number is sensitive for a general profile endpoint
			// "date_of_birth": user.DateOfBirth, // Consider if DOB is sensitive
		},
	})
}
