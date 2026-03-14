//go:build integration

package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ziqunw44/wechat/internal/model"
	"github.com/ziqunw44/wechat/internal/pool"
)

// seedGroupData 往 MySQL 写入测试用户、群组和群成员；若该用户已存在（上次运行残留）则复用，避免重复运行测试时 username 唯一约束冲突。
func seedGroupData(t *testing.T) (userUUID string, groupUUID string) {
	db := pool.GetDB()
	_ = db.AutoMigrate(&model.User{}, &model.Group{}, &model.GroupMember{})

	userUUID = "test-group-user-" + t.Name()
	username := "group_test_user_" + t.Name()

	var u model.User
	err := db.Where("uuid = ?", userUUID).First(&u).Error
	if err == nil {
		// 用户已存在，查其任意一个群
		var g model.Group
		if err := db.Table("groups").Joins("JOIN group_members ON group_members.group_id = groups.id").
			Where("group_members.user_id = ?", u.Id).Limit(1).First(&g).Error; err == nil {
			return userUUID, g.Uuid
		}
		// 有用户但没群，下面用 u.Id 补建群和成员
	} else {
		u = model.User{
			Uuid:     userUUID,
			Username: username,
			Password: "pass123",
			Nickname: "测试用户",
			CreateAt: time.Now(),
		}
		if err := db.Create(&u).Error; err != nil {
			t.Fatalf("seed user: %v", err)
		}
	}

	g := &model.Group{
		Uuid:      uuid.New().String(),
		UserId:    u.Id,
		Name:      "测试群",
		Notice:    "群公告",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(g).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	groupUUID = g.Uuid

	gm := &model.GroupMember{
		UserId:    u.Id,
		GroupId:   g.ID,
		Nickname:  u.Username,
		Mute:      0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(gm).Error; err != nil {
		t.Fatalf("seed group_member: %v", err)
	}
	return userUUID, groupUUID
}

func TestSeedGroupData(t *testing.T) {
	db := pool.GetDB()
	if db == nil {
		t.Fatal("DB nil，请先启动 MySQL 并配置 config.toml")
	}
	userUUID, groupUUID := seedGroupData(t)
	if userUUID == "" || groupUUID == "" {
		t.Fatal("seed 未返回 uuid")
	}
	t.Logf("已写入测试数据: userUUID=%s, groupUUID=%s", userUUID, groupUUID)
}

func TestGroupService_GetGroups(t *testing.T) {
	db := pool.GetDB()
	if db == nil {
		t.Fatal("DB nil")
	}
	userUUID, _ := seedGroupData(t)

	groups, err := GroupService.GetGroups(userUUID)
	if err != nil {
		t.Fatalf("GetGroups: %v", err)
	}
	if len(groups) == 0 {
		t.Fatal("GetGroups 返回 0 条，期望至少 1 条")
	}
	t.Logf("GetGroups 返回 %d 个群: %+v", len(groups), groups[0])
}

func TestGroupService_GetGroups_UserNotExist(t *testing.T) {
	_, err := GroupService.GetGroups("non-exist-uuid-xxx")
	if err == nil {
		t.Fatal("期望用户不存在时返回错误")
	}
	t.Logf("用户不存在时错误: %v", err)
}

func TestGroupService_SaveGroup(t *testing.T) {
	userUUID, _ := seedGroupData(t)

	g := model.Group{
		Name:   "SaveGroup 创建的群",
		Notice: "公告",
	}
	GroupService.SaveGroup(userUUID, g)

	groups, err := GroupService.GetGroups(userUUID)
	if err != nil {
		t.Fatalf("GetGroups after SaveGroup: %v", err)
	}
	var found bool
	for _, gr := range groups {
		if gr.Name == g.Name {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("SaveGroup 后未在 GetGroups 中找到群 %q", g.Name)
	}
	t.Log("SaveGroup 与 GetGroups 通过")
}

func TestGroupService_GetUserIdByGroupUuid(t *testing.T) {
	userUUID, groupUUID := seedGroupData(t)
	if groupUUID == "" {
		db := pool.GetDB()
		var g model.Group
		db.Where("name = ?", "测试群").First(&g)
		groupUUID = g.Uuid
	}

	users := GroupService.GetUserIdByGroupUuid(groupUUID)
	if len(users) == 0 {
		t.Fatal("GetUserIdByGroupUuid 返回 0 人")
	}
	t.Logf("群成员数: %d, userUUID: %s", len(users), userUUID)
}

func TestGroupService_JoinGroup(t *testing.T) {
	_, groupUUID := seedGroupData(t)
	db := pool.GetDB()

	joinUUID := "join-user-" + t.Name()
	var joinUser model.User
	if err := db.Where("uuid = ?", joinUUID).First(&joinUser).Error; err != nil {
		joinUser = model.User{
			Uuid:     joinUUID,
			Username: "join_user_" + t.Name(),
			Password: "pass",
			Nickname: "加群用户",
			CreateAt: time.Now(),
		}
		if err := db.Create(&joinUser).Error; err != nil {
			t.Fatalf("创建加群用户: %v", err)
		}
	}

	err := GroupService.JoinGroup(groupUUID, joinUser.Uuid)
	if err != nil && err.Error() != "已经加入该群组" {
		t.Fatalf("JoinGroup: %v", err)
	}

	users := GroupService.GetUserIdByGroupUuid(groupUUID)
	if len(users) < 2 {
		t.Fatalf("加群后群成员应至少 2 人，got %d", len(users))
	}
	t.Logf("JoinGroup 成功，当前群成员数: %d", len(users))

	// 重复加群应报错
	err = GroupService.JoinGroup(groupUUID, joinUser.Uuid)
	if err == nil {
		t.Fatal("重复加群应返回错误")
	}
	t.Logf("重复加群正确返回错误: %v", err)
}
