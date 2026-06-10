package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type qywxAppTokenResponse struct {
	Code     int    `json:"errcode"`
	Msg      string `json:"errmsg"`
	Token    string `json:"access_token"`
	ExpireIn int    `json:"expires_in"`
}

type qywxAppSendResponse struct {
	Code    int    `json:"errcode"`
	Msg     string `json:"errmsg"`
	MsgId   string `json:"msgid"`
}

type QyWeiXinApp struct {
	CorpID     string
	CorpSecret string
	AgentID    int
}

var (
	qywxAppTokenCache     = make(map[string]*cachedToken)
	qywxAppTokenCacheLock sync.Mutex
)

type cachedToken struct {
	token     string
	expiresAt time.Time
}

// getAccessToken 获取企业微信应用的access_token，带缓存
func (t *QyWeiXinApp) getAccessToken() (string, error) {
	cacheKey := t.CorpID + ":" + t.CorpSecret

	qywxAppTokenCacheLock.Lock()
	if cache, ok := qywxAppTokenCache[cacheKey]; ok && time.Now().Before(cache.expiresAt) {
		qywxAppTokenCacheLock.Unlock()
		return cache.token, nil
	}
	qywxAppTokenCacheLock.Unlock()

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", t.CorpID, t.CorpSecret)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var r qywxAppTokenResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", err
	}
	if r.Code != 0 {
		return "", fmt.Errorf("获取access_token失败: %s", string(body))
	}

	qywxAppTokenCacheLock.Lock()
	qywxAppTokenCache[cacheKey] = &cachedToken{
		token:     r.Token,
		expiresAt: time.Now().Add(time.Duration(r.ExpireIn-60) * time.Second),
	}
	qywxAppTokenCacheLock.Unlock()

	return r.Token, nil
}

func (t *QyWeiXinApp) send(body interface{}) ([]byte, error) {
	token, err := t.getAccessToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token)
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r qywxAppSendResponse
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return respBody, err
	}
	if r.Code != 0 {
		return respBody, fmt.Errorf("response error: %s", string(respBody))
	}
	return respBody, nil
}

// SendMessageText 发送文本消息
func (t *QyWeiXinApp) SendMessageText(text string, touser string, toparty string, totag string) ([]byte, error) {
	msg := map[string]interface{}{
		"touser":   touser,
		"toparty":  toparty,
		"totag":    totag,
		"msgtype":  "text",
		"agentid":  t.AgentID,
		"text":    map[string]interface{}{"content": text},
		"safe":    0,
	}
	return t.send(msg)
}

// SendMessageMarkdown 发送Markdown消息（个人）
func (t *QyWeiXinApp) SendMessageMarkdown(text string, touser string, toparty string, totag string) ([]byte, error) {
	msg := map[string]interface{}{
		"touser":   touser,
		"toparty":  toparty,
		"totag":    totag,
		"msgtype":  "markdown",
		"agentid":  t.AgentID,
		"markdown": map[string]interface{}{"content": text},
		"safe":    0,
	}
	return t.send(msg)
}

// CreateGroupChat 创建应用群聊
func (t *QyWeiXinApp) CreateGroupChat(name, owner string, userlist []string) (string, error) {
	token, err := t.getAccessToken()
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/appchat/create?access_token=%s", token)
	body := map[string]interface{}{
		"name":     name,
		"owner":    owner,
		"userlist": userlist,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var r struct {
		Code   int    `json:"errcode"`
		Msg    string `json:"errmsg"`
		ChatID string `json:"chatid"`
	}
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return "", fmt.Errorf("解析创建群聊响应失败: %s", string(respBody))
	}
	if r.Code != 0 {
		return "", fmt.Errorf("创建群聊失败: %s", string(respBody))
	}
	return r.ChatID, nil
}

// GetGroupChatInfo 获取群聊信息
func (t *QyWeiXinApp) GetGroupChatInfo(chatid string) (string, string, string, []string, error) {
	token, err := t.getAccessToken()
	if err != nil {
		return "", "", "", nil, err
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/appchat/get?access_token=%s&chatid=%s", token, chatid)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", "", nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", nil, err
	}

	var r struct {
		Code     int    `json:"errcode"`
		Msg      string `json:"errmsg"`
		ChatInfo struct {
			ChatID   string   `json:"chatid"`
			Name     string   `json:"name"`
			Owner    string   `json:"owner"`
			Userlist []string `json:"userlist"`
		} `json:"chat_info"`
	}
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return "", "", "", nil, fmt.Errorf("解析群聊信息失败: %s", string(respBody))
	}
	if r.Code != 0 {
		return "", "", "", nil, fmt.Errorf("获取群聊信息失败: %s", string(respBody))
	}
	return r.ChatInfo.ChatID, r.ChatInfo.Name, r.ChatInfo.Owner, r.ChatInfo.Userlist, nil
}

// sendToGroup 向群聊发送消息
func (t *QyWeiXinApp) sendToGroup(body interface{}) ([]byte, error) {
	token, err := t.getAccessToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/appchat/send?access_token=%s", token)
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r qywxAppSendResponse
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return respBody, err
	}
	if r.Code != 0 {
		return respBody, fmt.Errorf("response error: %s", string(respBody))
	}
	return respBody, nil
}

// SendGroupText 向群聊发送文本消息
func (t *QyWeiXinApp) SendGroupText(text string, chatid string) ([]byte, error) {
	msg := map[string]interface{}{
		"chatid":  chatid,
		"msgtype": "text",
		"text":    map[string]interface{}{"content": text},
		"safe":    0,
	}
	return t.sendToGroup(msg)
}

// SendGroupMarkdown 向群聊发送Markdown消息
func (t *QyWeiXinApp) SendGroupMarkdown(text string, chatid string) ([]byte, error) {
	msg := map[string]interface{}{
		"chatid":   chatid,
		"msgtype":  "markdown",
		"markdown": map[string]interface{}{"content": text},
		"safe":     0,
	}
	return t.sendToGroup(msg)
}
