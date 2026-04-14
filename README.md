# wechat

基于 Go（Gin、GORM 等）与 Vue 的即时通讯相关服务端与 Web 前端示例项目。

## 运行后端

1. 安装并启动 **MySQL**，按 `configs/config.toml` 创建库并填写账号密码。  
2. 在项目根目录执行：

```bash
go run .
```

默认监听配置中的端口（如 **8000**）。`GET /` 返回简单健康信息；登录、会话、WebSocket 等接口见 `internal/https_server/https_server.go`。

开发环境默认 **Gin debug**，不会把 HTTP 强制跳转到 HTTPS，便于使用 `http://127.0.0.1:8000` 调试。上线可设置 `export GIN_MODE=release`。

## 联调前端

修改 `web/chat-server/src/store/index.js` 中的 `backendUrl`、`wsUrl`，使其指向本机后端（例如 `http://127.0.0.1:8000` 与 `ws://127.0.0.1:8000`）。前端 `vue.config.js` 若使用系统证书路径，本地开发可按注释改为常见端口与证书配置。

## 许可证

SPDX-License-Identifier: GPL-3.0-or-later

本程序为自由软件：你可以在遵循许可证的前提下再发布和/或修改它。

```
Copyright (C) 2026  wangziqun

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```

完整许可证条文见仓库根目录的 [LICENSE](LICENSE) 文件。
