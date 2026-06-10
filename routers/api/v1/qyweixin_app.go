package v1

import (
	"encoding/json"
	"message-nest/models"
	"message-nest/pkg/app"
	"message-nest/pkg/message"
	"message-nest/service/send_way_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateGroupChat 创建企业微信应用群聊
func CreateGroupChat(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		WayID    string   `json:"way_id" validate:"required,max=12"`
		Name     string   `json:"name" validate:"required,max=50"`
		Owner    string   `json:"owner" validate:"required,max=100"`
		Userlist []string `json:"userlist" validate:"required,min=2"`
	}
	errCode, errMsg := app.BindJsonAndPlayValid(c, &req)
	if errCode != 200 {
		appG.CResponse(errCode, errMsg, nil)
		return
	}

	// 获取渠道信息
	wayObj, err := models.GetWayByID(req.WayID)
	if err != nil {
		appG.CResponse(http.StatusInternalServerError, "获取渠道信息失败", nil)
		return
	}

	var auth send_way_service.WayDetailQyWeiXinApp
	err = json.Unmarshal([]byte(wayObj.Auth), &auth)
	if err != nil {
		appG.CResponse(http.StatusBadRequest, "渠道认证信息解析失败", nil)
		return
	}

	cli := message.QyWeiXinApp{
		CorpID:     auth.CorpID,
		CorpSecret: auth.CorpSecret,
		AgentID:    auth.AgentID,
	}

	chatID, err := cli.CreateGroupChat(req.Name, req.Owner, req.Userlist)
	if err != nil {
		appG.CResponse(http.StatusInternalServerError, err.Error(), nil)
		return
	}

	userlistJSON, _ := json.Marshal(req.Userlist)
	err = models.AddQyWeiXinAppChat(req.WayID, chatID, req.Name, req.Owner, string(userlistJSON))
	if err != nil {
		appG.CResponse(http.StatusInternalServerError, "保存群聊信息失败", nil)
		return
	}

	appG.CResponse(http.StatusOK, "创建群聊成功", map[string]string{"chatid": chatID, "name": req.Name})
}

// ListGroupChats 获取群聊列表
func ListGroupChats(c *gin.Context) {
	appG := app.Gin{C: c}
	wayID := c.Query("way_id")

	if wayID == "" {
		appG.CResponse(http.StatusBadRequest, "缺少way_id参数", nil)
		return
	}

	chats, err := models.GetQyWeiXinAppChatsByWayID(wayID)
	if err != nil {
		appG.CResponse(http.StatusInternalServerError, "获取群聊列表失败", nil)
		return
	}

	if chats == nil {
		chats = []models.QyWeiXinAppChat{}
	}

	appG.CResponse(http.StatusOK, "获取群聊列表成功", chats)
}

// RefreshGroupChat 刷新群聊信息
func RefreshGroupChat(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		ID    int    `json:"id" validate:"required"`
		WayID string `json:"way_id" validate:"required,max=12"`
	}
	errCode, errMsg := app.BindJsonAndPlayValid(c, &req)
	if errCode != 200 {
		appG.CResponse(errCode, errMsg, nil)
		return
	}

	chat, err := models.GetQyWeiXinAppChatByID(req.ID)
	if err != nil {
		appG.CResponse(http.StatusInternalServerError, "群聊不存在", nil)
		return
	}

	wayObj, err := models.GetWayByID(req.WayID)
	if err != nil {
		appG.CResponse(http.StatusInternalServerError, "获取渠道信息失败", nil)
		return
	}

	var auth send_way_service.WayDetailQyWeiXinApp
	err = json.Unmarshal([]byte(wayObj.Auth), &auth)
	if err != nil {
		appG.CResponse(http.StatusBadRequest, "渠道认证信息解析失败", nil)
		return
	}

	cli := message.QyWeiXinApp{
		CorpID:     auth.CorpID,
		CorpSecret: auth.CorpSecret,
		AgentID:    auth.AgentID,
	}

	_, name, owner, userlist, err := cli.GetGroupChatInfo(chat.ChatID)
	if err != nil {
		appG.CResponse(http.StatusInternalServerError, err.Error(), nil)
		return
	}

	userlistJSON, _ := json.Marshal(userlist)
	data := map[string]interface{}{
		"name":     name,
		"owner":    owner,
		"userlist": string(userlistJSON),
	}
	err = models.UpdateQyWeiXinAppChat(req.ID, data)
	if err != nil {
		appG.CResponse(http.StatusInternalServerError, "更新群聊信息失败", nil)
		return
	}

	chat.Name = name
	chat.Owner = owner
	chat.Userlist = string(userlistJSON)

	appG.CResponse(http.StatusOK, "刷新群聊成功", chat)
}

// DeleteGroupChat 删除群聊记录
func DeleteGroupChat(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		ID int `json:"id" validate:"required"`
	}
	errCode, errMsg := app.BindJsonAndPlayValid(c, &req)
	if errCode != 200 {
		appG.CResponse(errCode, errMsg, nil)
		return
	}

	err := models.DeleteQyWeiXinAppChat(req.ID)
	if err != nil {
		appG.CResponse(http.StatusInternalServerError, "删除群聊失败", nil)
		return
	}

	appG.CResponse(http.StatusOK, "删除群聊成功", nil)
}
