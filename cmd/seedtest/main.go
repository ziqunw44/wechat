// 插入或重置本地测试账号（密码明文，与 Login 逻辑一致）。
//
// 用法：在项目根目录执行 go run ./cmd/seedtest
//
// Copyright (C) 2026  wangziqun
// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"wechat/internal/config"
	"wechat/internal/dao"
	"wechat/internal/model"
	"wechat/pkg/enum/user_info/user_status_enum"
)

const (
	testTelephone = "13800138000"
	testPassword  = "test123456"
	testNickname  = "测试账号"
	testUUID      = "Utest00000001"
)

func main() {
	config.LoadConfig()
	if err := dao.InitGorm(); err != nil {
		panic(err)
	}

	var user model.UserInfo
	err := dao.GormDB.Where("telephone = ?", testTelephone).First(&user).Error
	if err == nil {
		user.Password = testPassword
		user.Nickname = testNickname
		if err := dao.GormDB.Save(&user).Error; err != nil {
			panic(err)
		}
		fmt.Printf("已存在手机号 %s，已重置密码。\n登录：手机 %s  密码 %s\n", testTelephone, testTelephone, testPassword)
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	}

	user = model.UserInfo{
		Uuid:      testUUID,
		Nickname:  testNickname,
		Telephone: testTelephone,
		Password:  testPassword,
		Avatar:    "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png",
		CreatedAt: time.Now(),
		IsAdmin:   0,
		Status:    user_status_enum.NORMAL,
	}
	if err := dao.GormDB.Create(&user).Error; err != nil {
		panic(err)
	}
	fmt.Printf("已创建测试账号。\n登录：手机 %s  密码 %s\n", testTelephone, testPassword)
}
