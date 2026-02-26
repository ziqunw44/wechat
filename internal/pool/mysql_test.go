//go:build integration

package pool

import (
	"testing"
	"time"

	"github.com/ziqunw44/wechat/internal/model"
)

// 集成测试：需要本地 MySQL 已启动，且 config.toml 中配置正确。
//
// 1. 确保 MySQL 运行，并已创建数据库：
//
//	mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS chat;"
//
// 2. 修改 config.toml 中的 [mysql] 与本地账号密码一致（如 user、password、port）。
//
// 3. 在项目根目录执行：
//
//	go test -tags=integration ./internal/pool/ -v
//
// 不带 -tags=integration 时不会编译本文件，不会连接数据库。
func TestGetDB(t *testing.T) {
	db := GetDB()
	if db == nil {
		t.Fatal("GetDB() 返回 nil")
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() 失败: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("数据库 Ping 失败: %v", err)
	}
	t.Log("数据库连接正常")
}

// TestAutoMigrate 检查 GORM 与当前模型的表结构是否一致（不修改数据）
func TestAutoMigrate(t *testing.T) {
	db := GetDB()
	err := db.AutoMigrate(
		&model.User{},
		&model.UserFriend{},
		&model.Message{},
		&model.Group{},
		&model.GroupMember{},
	)
	if err != nil {
		t.Fatalf("AutoMigrate 失败: %v", err)
	}
	t.Log("AutoMigrate 成功")
}

// TestUserCRUD 简单验证 User 表的增删查（写入后查询再删除）
func TestUserCRUD(t *testing.T) {
	db := GetDB()
	u := &model.User{
		Uuid:     "test-db-uuid-" + t.Name(),
		Username: "test_user_crud_" + t.Name(),
		Password: "testpass",
		Avatar:   "",
		CreateAt: time.Now(),
	}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("Create User 失败: %v", err)
	}
	t.Logf("创建用户 id=%d", u.Id)

	var got model.User
	if err := db.Where("id = ?", u.Id).First(&got).Error; err != nil {
		t.Fatalf("查询 User 失败: %v", err)
	}
	if got.Username != u.Username {
		t.Fatalf("查询结果 username 不一致: got %q", got.Username)
	}

	if err := db.Delete(&model.User{}, u.Id).Error; err != nil {
		t.Logf("清理测试用户失败（可忽略）: %v", err)
	}
	t.Log("User CRUD 通过")
}
