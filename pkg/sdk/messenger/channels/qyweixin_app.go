package channels

import (
	"fmt"
	"message-nest/pkg/message"
)

type QyWeiXinAppChannel struct{ *BaseChannel }

func NewQyWeiXinAppChannel() Channel {
	return &QyWeiXinAppChannel{NewBaseChannel(ChannelQyWeiXinApp, []string{FormatTypeMarkdown, FormatTypeText})}
}

func (c *QyWeiXinAppChannel) Send(config ChannelConfig, msg *Message) (*Result, error) {
	corpID := config.GetString("corpid")
	corpSecret := config.GetString("corpsecret")
	agentIDStr := config.GetString("agentid")

	if corpID == "" {
		return SendError("qyweixinapp config missing: corpid is required"), nil
	}
	if corpSecret == "" {
		return SendError("qyweixinapp config missing: corpsecret is required"), nil
	}
	if agentIDStr == "" {
		return SendError("qyweixinapp config missing: agentid is required"), nil
	}

	var agentID int
	_, err := fmt.Sscanf(agentIDStr, "%d", &agentID)
	if err != nil {
		return SendError("qyweixinapp config: agentid must be a number"), nil
	}

	touser := config.GetString("touser")
	toparty := config.GetString("toparty")
	totag := config.GetString("totag")

	contentType, formattedContent := c.FormatContent(msg)

	cli := message.QyWeiXinApp{
		CorpID:     corpID,
		CorpSecret: corpSecret,
		AgentID:    agentID,
	}
	var res []byte

	if contentType == FormatTypeText {
		res, err = cli.SendMessageText(formattedContent, touser, toparty, totag)
	} else if contentType == FormatTypeMarkdown {
		res, err = cli.SendMessageMarkdown(formattedContent, touser, toparty, totag)
	} else {
		return SendError("未知的企业微信应用发送内容类型：%s", contentType), nil
	}

	if err != nil {
		return ErrorResult(string(res), err), nil
	}
	return SuccessResult(string(res)), nil
}
