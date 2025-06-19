package models

import "gorm.io/gorm"

// User ...
type User struct {
	Username      string `json:"username" gorm:"primary_key"`
	Password      string `json:"password"`
	Prefix        string `json:"prefix"`
	First_name    string `json:"first_name"`
	Last_name     string `json:"last_name"`
	Email         string `json:"email"`
	Phone_number  string `json:"phone_number"`
	Date_of_birth string `json:"date_of_birth"`
}

func (u *User) TableName() string {
	return "users"
}

// GetAllUsers Fetch all User data
func GetAllUsers(db *gorm.DB, users *[]User) (err error) {
	if err := db.Raw("SELECT * FROM users").Scan(users).Error; err != nil {
		return err
	}
	return nil
}

func CreateUser(db *gorm.DB, users *User) (err error) { // Insert New data
	if err := db.Create(&users).Error; err != nil {
		return err
	}
	return nil
}

// GetUserByID ... Fetch only one user by Id
func GetUserByID(db *gorm.DB, users *User, username string) error { // Use parameterized query
	if err := db.Raw("SELECT * FROM users WHERE username = ?", username).First(&users).Error; err != nil {
		return err // Return the original GORM error
	}
	return nil
}

func UpdateUser(db *gorm.DB, user *User) (err error) { // Update user (Assuming user.Username is the primary key)
	// ทำการบันทึกข้อมูลผู้ใช้
	if err := db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUser(db *gorm.DB, username string) (err error) { // Delete user (Assuming username is the primary key)
	return db.Where("username = ?", username).Delete(&User{}).Error
}
