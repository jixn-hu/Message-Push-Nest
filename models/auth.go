package models

import (
	"gorm.io/gorm"
)

type Auth struct {
	ID       int    `json:"id" gorm:"autoIncrement;type:integer ;primaryKey" json:"id"`
	Username string `json:"username" gorm:"type:varchar(100);default:'';"`
	Password string `json:"password" gorm:"type:varchar(100);default:'';"`
	Role     string `json:"role" gorm:"type:varchar(20);default:'admin'"`
}

// CheckAuth 检查用户信息
func CheckAuth(username, password string) (bool, error) {
	var auth Auth
	err := db.Select("id").Where(Auth{Username: username, Password: password}).Limit(1).Find(&auth).Error
	if err != nil {
		return false, err
	}

	if auth.ID > 0 {
		return true, nil
	}

	return false, nil
}

func GetUserByUsername(username string) (*Auth, error) {
    var auth Auth
    err := db.Where("username = ?", username).Limit(1).Find(&auth).Error
    if err != nil {
        return nil, err
    }
	if auth.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
    return &auth, nil
}

// EditUser 编辑用户信息
func EditUser(username string, data map[string]interface{}) error {
	if err := db.Model(&Auth{}).Where("username = ? ", username).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func AddUser(account string, password string) error {
	auth := Auth{
		Username: account,
		Password: password,
	}
	if err := db.Create(&auth).Error; err != nil {
		return err
	}
	return nil
}

func AddUserWithRole(account, password, role string) error {
	auth := Auth{
		Username: account,
		Password: password,
		Role:     role,
	}
	if err := db.Create(&auth).Error; err != nil {
		return err
	}
	return nil
}

func GetAllUsers() ([]Auth, error) {
	var users []Auth
	err := db.Order("id ASC").Find(&users).Error
	return users, err
}

func DeleteUserByID(id int) error {
	return db.Where("id = ?", id).Delete(&Auth{}).Error
}
