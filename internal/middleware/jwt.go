package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"order-notification-system/internal/auth" // Import the new auth package

	"github.com/gin-gonic/gin"
)

// JWTMiddleware validates JWT token in request header
func JWTMiddleware() gin.HandlerFunc {
	// Add Line Numbers to Log Output
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return func(c *gin.Context) {
		var tokenValue string

		// ตรวจสอบว่าเป็น WebSocket upgrade request หรือไม่
		isWebSocketUpgrade := c.GetHeader("Upgrade") == "websocket" &&
			strings.Contains(strings.ToLower(c.GetHeader("Connection")), "upgrade")

		if isWebSocketUpgrade {
			tokenValue = c.Query("token") // ดึง token จาก query parameter ชื่อ "token"
			if tokenValue == "" {
				// สำหรับ WebSocket ถ้าไม่มี token ใน query parameter ให้ปฏิเสธการเชื่อมต่อ
				// Client ควรส่ง token มาในรูปแบบ ws://host/ws?token=YOUR_TOKEN
				// การตอบกลับด้วย JSON อาจจะไม่ถูกเห็นโดย client ws.onerror โดยตรง
				// แต่การเชื่อมต่อจะล้มเหลว และ server ควร abort handshake
				// log.Println("WebSocket connection attempt without token in query.")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status":  "error",
					"message": "WebSocket authentication token is required as a query parameter 'token'",
				})
				return
			}
		} else {
			// สำหรับ HTTP request ปกติ ดึง token จาก Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  "error",
					"message": "Authorization header is required and must be Bearer token",
				})
				c.Abort()
				return
			}
			tokenValue = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if tokenValue == "" { // ควรถูกดักจับโดยเงื่อนไขด้านบนแล้ว แต่เป็นการป้องกันอีกชั้น
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Authentication token not provided or improperly formatted",
			})
			c.Abort()
			return
		}

		// Use the new auth package for verification
		if os.Getenv("JWT_SECRET_KEY") == "" {
			log.Println("CRITICAL: JWT_SECRET_KEY is not set in environment variables.")
		}
		token, claims, err := auth.VerifyToken(tokenValue)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Store claims in context for handlers to use
		// claims is now *auth.CustomClaims from auth.VerifyToken
		c.Set("claims", claims)
		c.Next()
	}
}
