package models

import "message-nest/pkg/util"

type QyWeiXinAppChat struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	WayID    string `json:"way_id" gorm:"type:varchar(12);default:'';index"`
	ChatID   string `json:"chatid" gorm:"type:varchar(100);default:'';uniqueIndex"`
	Name     string `json:"name" gorm:"type:varchar(200);default:''"`
	Owner    string `json:"owner" gorm:"type:varchar(100);default:''"`
	Userlist string `json:"userlist" gorm:"type:text"`

	CreatedAt util.Time `json:"created_on" gorm:"column:created_on;autoCreateTime"`
	UpdatedAt util.Time `json:"modified_on" gorm:"column:modified_on;autoUpdateTime"`
}

func AddQyWeiXinAppChat(wayID, chatID, name, owner, userlist string) error {
	chat := QyWeiXinAppChat{
		WayID:    wayID,
		ChatID:   chatID,
		Name:     name,
		Owner:    owner,
		Userlist: userlist,
	}
	return db.Create(&chat).Error
}

func GetQyWeiXinAppChatsByWayID(wayID string) ([]QyWeiXinAppChat, error) {
	var chats []QyWeiXinAppChat
	err := db.Where("way_id = ?", wayID).Order("created_on DESC").Find(&chats).Error
	return chats, err
}

func GetQyWeiXinAppChatByID(id int) (*QyWeiXinAppChat, error) {
	var chat QyWeiXinAppChat
	err := db.Where("id = ?", id).First(&chat).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func UpdateQyWeiXinAppChat(id int, data map[string]interface{}) error {
	return db.Model(&QyWeiXinAppChat{}).Where("id = ?", id).Updates(data).Error
}

func DeleteQyWeiXinAppChat(id int) error {
	return db.Where("id = ?", id).Delete(&QyWeiXinAppChat{}).Error
}
