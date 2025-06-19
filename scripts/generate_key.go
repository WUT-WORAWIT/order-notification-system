package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

func generateSecureKey(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func main() {
	key := generateSecureKey(32)

	envContent := fmt.Sprintf(`# Security
## JWT Configuration
JWT_SECRET_KEY="%s"    # Do not regenerate this key
JWT_EXPIRATION=24h     # Token expiration time`, key)

	if err := os.WriteFile(".env.key", []byte(envContent), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generated new JWT secret key in .env.key")
}

// กรณีที่ควรเปลี่ยน JWT Secret Key
// 1. กรณีด้านความปลอดภัย
// เมื่อมีการรั่วไหลของ Secret Key
// เมื่อพนักงานที่รู้ Key ลาออก
// เมื่อพบว่ามีการใช้ Token ในทางที่ผิด
// เมื่อระบบถูกโจมตี
// 2. กรณีเปลี่ยนสภาพแวดล้อม
// เมื่อ Deploy ขึ้น Production ครั้งแรก
// เมื่อย้ายไป Server ใหม่
// เมื่อเปลี่ยน Environment (dev → staging → prod)
// 3. กรณีตามนโยบาย
// ตามรอบการเปลี่ยน Key ที่กำหนด (เช่น ทุก 6 เดือน)
// เมื่อมีการ Audit ความปลอดภัย
// เมื่อต้องทำตามมาตรฐานความปลอดภัยใหม่
