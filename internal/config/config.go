package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Init initializes database connection
func Init() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// เก็บ DSN ทั้งสองแบบไว้ใน environment variables
	localDSN := os.Getenv("LOCAL_DSN")
	dockerDSN := os.Getenv("DOCKER_DSN")

	// เลือก DSN ตาม environment
	var dsn string
	if os.Getenv("APP_ENV") == "docker" {
		dsn = dockerDSN // "host=host.docker.internal user=postgres password=123456 dbname=TEST port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	} else {
		dsn = localDSN // "host=127.0.0.1 user=postgres password=123456 dbname=TEST port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	}

	if dsn == "" {
		log.Fatal("❌ DSN is empty. Please check your environment variables.")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// เพิ่ม configuration options ตามต้องการ
	})

	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// ทดสอบการเชื่อมต่อ
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get database instance: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	fmt.Println("✅ Connected to database successfully")

	return db
}

// Close closes database connection
func Close(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("❌ Error getting database instance: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("❌ Error closing database connection: %v", err)
		return
	}
	fmt.Println("✅ Database connection closed successfully")
}
