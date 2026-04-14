// 批量写入演示用户与群聊（明文密码，与 Login 一致）。
//
// 手机号需满足前端校验 ^1[3456789]\d{9}$，故使用 13100000121～13100000124，
// 密码均为 123。用法：go run ./cmd/seedfixture
//
// Copyright (C) 2026  wangziqun
// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"wechat/internal/config"
	"wechat/internal/dao"
	"wechat/internal/model"
	"wechat/pkg/enum/contact/contact_status_enum"
	"wechat/pkg/enum/contact/contact_type_enum"
	"wechat/pkg/enum/group_info/add_mode_enum"
	"wechat/pkg/enum/group_info/group_status_enum"
	"wechat/pkg/enum/user_info/user_status_enum"
)

const defaultAvatar = "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"

type demoUser struct {
	UUID      string
	Telephone string
	Password  string
	Nickname  string
}

var demoUsers = []demoUser{
	{"Udemo00000001", "13100000121", "123", "用户甲"},
	{"Udemo00000002", "13100000122", "123", "用户乙"},
	{"Udemo00000003", "13100000123", "123", "用户丙"},
	{"Udemo00000004", "13100000124", "123", "用户丁"},
}

type demoGroup struct {
	UUID    string
	Name    string
	Notice  string
	OwnerID string
	Members []string // 含群主，顺序不限，会去重
}

var demoGroups = []demoGroup{
	{"Gdemo00000001", "同学群", "测试群公告：同学交流", "Udemo00000001", []string{"Udemo00000001", "Udemo00000002", "Udemo00000003"}},
	{"Gdemo00000002", "项目一组", "联调用项目群", "Udemo00000002", []string{"Udemo00000002", "Udemo00000001", "Udemo00000004"}},
	{"Gdemo00000003", "闲聊灌水", "随便聊聊", "Udemo00000003", []string{"Udemo00000003", "Udemo00000001", "Udemo00000002", "Udemo00000004"}},
}

func main() {
	config.LoadConfig()
	if err := dao.InitGorm(); err != nil {
		panic(err)
	}

	for _, u := range demoUsers {
		if err := upsertUser(u); err != nil {
			panic(err)
		}
	}

	for _, g := range demoGroups {
		if err := upsertGroup(g); err != nil {
			panic(err)
		}
	}

	fmt.Println("--- 演示数据已就绪 ---")
	fmt.Println("密码均为: 123")
	for _, u := range demoUsers {
		fmt.Printf("  %s  | 手机 %s  | %s\n", u.Nickname, u.Telephone, u.UUID)
	}
	fmt.Println("群聊:")
	for _, g := range demoGroups {
		fmt.Printf("  [%s] %s  owner=%s  成员数=%d\n", g.UUID, g.Name, g.OwnerID, len(uniqueStrings(g.Members)))
	}
}

func upsertUser(u demoUser) error {
	var row model.UserInfo
	err := dao.GormDB.Where("telephone = ?", u.Telephone).First(&row).Error
	if err == nil {
		row.Password = u.Password
		row.Nickname = u.Nickname
		row.Status = user_status_enum.NORMAL
		return dao.GormDB.Save(&row).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	err = dao.GormDB.Where("uuid = ?", u.UUID).First(&row).Error
	if err == nil {
		row.Telephone = u.Telephone
		row.Password = u.Password
		row.Nickname = u.Nickname
		row.Status = user_status_enum.NORMAL
		return dao.GormDB.Save(&row).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	row = model.UserInfo{
		Uuid:      u.UUID,
		Nickname:  u.Nickname,
		Telephone: u.Telephone,
		Password:  u.Password,
		Avatar:    defaultAvatar,
		CreatedAt: time.Now(),
		IsAdmin:   0,
		Status:    user_status_enum.NORMAL,
	}
	return dao.GormDB.Create(&row).Error
}

func upsertGroup(g demoGroup) error {
	members := uniqueStrings(g.Members)
	if len(members) == 0 {
		return fmt.Errorf("群 %s 无成员", g.UUID)
	}
	raw, err := json.Marshal(members)
	if err != nil {
		return err
	}

	var group model.GroupInfo
	err = dao.GormDB.Where("uuid = ?", g.UUID).First(&group).Error
	now := time.Now()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		group = model.GroupInfo{
			Uuid:      g.UUID,
			Name:      g.Name,
			Notice:    g.Notice,
			Members:   raw,
			MemberCnt: len(members),
			OwnerId:   g.OwnerID,
			AddMode:   add_mode_enum.DIRECT,
			Avatar:    defaultAvatar,
			Status:    group_status_enum.NORMAL,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := dao.GormDB.Create(&group).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		group.Name = g.Name
		group.Notice = g.Notice
		group.Members = raw
		group.MemberCnt = len(members)
		group.OwnerId = g.OwnerID
		group.AddMode = add_mode_enum.DIRECT
		group.Status = group_status_enum.NORMAL
		group.UpdatedAt = now
		if err := dao.GormDB.Save(&group).Error; err != nil {
			return err
		}
	}

	for _, uid := range members {
		if err := ensureGroupContact(uid, g.UUID); err != nil {
			return err
		}
	}
	return nil
}

func ensureGroupContact(userID, groupUUID string) error {
	var c model.UserContact
	err := dao.GormDB.Where("user_id = ? AND contact_id = ? AND contact_type = ?", userID, groupUUID, contact_type_enum.GROUP).
		First(&c).Error
	if err == nil {
		if c.Status != contact_status_enum.NORMAL {
			c.Status = contact_status_enum.NORMAL
			c.UpdateAt = time.Now()
			return dao.GormDB.Save(&c).Error
		}
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	now := time.Now()
	c = model.UserContact{
		UserId:      userID,
		ContactId:   groupUUID,
		ContactType: contact_type_enum.GROUP,
		Status:      contact_status_enum.NORMAL,
		CreatedAt:   now,
		UpdateAt:    now,
	}
	return dao.GormDB.Create(&c).Error
}

func uniqueStrings(in []string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, s := range in {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
