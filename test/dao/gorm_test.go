package dao

import (
	"wechat/internal/dao"
	"wechat/internal/model"
	"wechat/pkg/util/random"
	"strconv"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	if err := dao.InitGorm(); err != nil {
		t.Skipf("mysql unavailable: %v", err)
	}
	userInfo := &model.UserInfo{
		Uuid:      "U" + strconv.Itoa(random.GetRandomInt(11)),
		Nickname:  "apylee",
		Telephone: "18032353211",
		Email:     "1212312312@qq.com",
		Password:  "123456",
		CreatedAt: time.Now(),
		IsAdmin:   1,
		Status:    0,
	}
	_ = dao.GormDB.Create(userInfo)
}
