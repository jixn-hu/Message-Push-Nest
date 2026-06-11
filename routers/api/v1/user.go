package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"message-nest/models"
	"message-nest/pkg/app"
)

// ListUsers 获取用户列表
func ListUsers(c *gin.Context) {
	appG := app.Gin{C: c}
	users, err := models.GetAllUsers()
	if err != nil {
		appG.CResponse(http.StatusInternalServerError, "获取用户列表失败", nil)
		return
	}
	if users == nil {
		users = []models.Auth{}
	}
	appG.CResponse(http.StatusOK, "获取用户列表成功", users)
}

// AddUserReq 新增用户请求
type AddUserReq struct {
	Username string `json:"username" validate:"required,max=50" label:"用户名"`
	Password string `json:"password" validate:"required,max=50" label:"密码"`
	Role     string `json:"role" validate:"required,oneof=admin user" label:"角色"`
}

// AddNewUser 新增用户
func AddNewUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var req AddUserReq
	errCode, errMsg := app.BindJsonAndPlayValid(c, &req)
	if errCode != 200 {
		appG.CResponse(errCode, errMsg, nil)
		return
	}

	err := models.AddUserWithRole(req.Username, req.Password, req.Role)
	if err != nil {
		appG.CResponse(http.StatusInternalServerError, "添加用户失败: "+err.Error(), nil)
		return
	}
	appG.CResponse(http.StatusOK, "添加用户成功", nil)
}

// DeleteUserReq 删除用户请求
type DeleteUserReq struct {
	ID int `json:"id" validate:"required" label:"用户ID"`
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var req DeleteUserReq
	errCode, errMsg := app.BindJsonAndPlayValid(c, &req)
	if errCode != 200 {
		appG.CResponse(errCode, errMsg, nil)
		return
	}

	err := models.DeleteUserByID(req.ID)
	if err != nil {
		appG.CResponse(http.StatusInternalServerError, "删除用户失败: "+err.Error(), nil)
		return
	}
	appG.CResponse(http.StatusOK, "删除用户成功", nil)
}
