# PushX

消息推送整合平台，统一管理和发送多种渠道的消息通知。

原名 Message Nest（信息巢），意为一个拥有各种渠道信息方式的集合站点。如果你有多种消息推送方式，每次都需要调用不同接口去发送消息，这个项目可以帮你管理各种消息方式，提供统一的发送 API，一个请求推送到多种渠道。

## 特色

- **多渠道整合**：支持邮件、钉钉、企业微信、企业微信应用、飞书、Telegram、Bark、PushMe、Gotify、Ntfy、阿里云短信等
- **消息模板**：通过占位符实现动态内容替换，一次定义多处复用
- **企业微信应用**：支持个人发送和群聊发送，支持新建群聊、同步群聊信息
- **角色权限**：管理员 / 普通用户，普通用户不可编辑渠道和系统设置
- **用户管理**：管理员可新增、删除用户并分配角色
- **易于扩展**：模块化设计，方便接入新的消息渠道

## 快速开始

### Docker 部署

```bash
docker run -d \
  --name pushx \
  -p 8394:8000 \
  -v ./data:/app/conf \
  crpi-qcglwoi73qr0h01l.cn-hangzhou.personal.cr.aliyuncs.com/jixn/pushx:latest
```

或使用 docker-compose：

```yaml
version: "3.7"
services:
  pushx:
    image: crpi-qcglwoi73qr0h01l.cn-hangzhou.personal.cr.aliyuncs.com/jixn/pushx:latest
    container_name: pushx
    restart: always
    ports:
      - "8394:8000"
    volumes:
      - ./data:/app/conf
```

默认账号：`admin` / `123456`

### 本地开发

```bash
# 后端
go run main.go

# 前端
cd web && npm install && npm run dev
```

## 更新日志

- **2026.06.09** - 品牌改名 PushX、新增企业微信应用渠道（群聊管理）、角色权限系统、用户管理、前端导航精简
- **2026.06.03** - 新增 Gotify/Ntfy 渠道、自定义主题色、推送 SDK API、Windows 字体优化
- **2026.02.04** - 新增 PushMe 推送渠道
- **2026.01.21** - 新增 Bark 推送渠道
- **2026.01.15** - 新增 URL 路径前缀配置
- **2026.01.14** - 新增 Telegram 机器人渠道、HTTP/SOCKS5 代理

[查看完整更新日志](docs/guide/changelog.md)

## 许可证

本项目基于 [Apache-2.0](LICENSE) 开源协议，且包含 [NOTICE](NOTICE) 文件。

**强制要求：** 任何形式的使用、二次开发、修改或分发，均**必须**在显著位置标注作者来源（engigu）并保留原始仓库链接。详情请参阅 [NOTICE](NOTICE) 文件。

## 项目来源

本项目基于 [engigu/Message-Push-Nest](https://github.com/engigu/Message-Push-Nest) 二次开发。
