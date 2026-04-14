// wechat — chat server and related services.
//
// Copyright (C) 2026  wangziqun
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see https://www.gnu.org/licenses/.
package main

import (
	"strconv"
	"wechat/internal/config"
	"wechat/internal/dao"
	"wechat/internal/https_server"
	"wechat/internal/service/chat"
	"wechat/internal/service/kafka"
)

func main() {
	// https_server.init 会调用 config.GetConfig() 并预载配置；此处再显式加载一次保证顺序清晰
	config.LoadConfig()
	dao.MustInitGorm()

	cfg := config.GetConfig()
	port := strconv.Itoa(cfg.MainConfig.Port)

	// 消费 WebSocket 投递的聊天消息并广播回显；未启动则仅能在 client.Read 里打到日志，不会入库、不会推回前端
	if cfg.KafkaConfig.MessageMode == "channel" {
		go chat.ChatServer.Start()
	} else {
		kafka.KafkaService.KafkaInit()
		go chat.KafkaChatServer.Start()
	}

	// 路由在 https_server.init 中注册到 GE。生产环境建议: export GIN_MODE=release（启用 HTTPS 跳转等安全策略）
	if err := https_server.GE.Run(":" + port); err != nil {
		panic(err)
	}
}
