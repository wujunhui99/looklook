package main

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/wujunhui99/looklook/app/usercenter/cmd/rpc/internal/config"
	"github.com/wujunhui99/looklook/app/usercenter/model"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 生成唯一的手机号
func generateUniqueMobile() string {
	return fmt.Sprintf("139%08d", rand.Int63n(100000000))
}

// 配置文件路径
const (
	testConfigFile = "etc/usercenter.yaml"
)

var (
	testUserModel     model.UserModel
	testUserAuthModel model.UserAuthModel
	testConfig        config.Config
)

// 初始化测试环境
func setupTestDB(t *testing.T) {
	// 读取配置文件
	var c config.Config
	err := conf.Load(testConfigFile, &c)
	if err != nil {
		t.Fatalf("加载配置文件失败: %v", err)
	}
	testConfig = c

	// 手动设置DSN（修复时区问题，使用localhost:33069）
	// MySQL在Docker中，端口映射到33069
	dsn := "root:PXDN93VRKUm8TeE7@tcp(localhost:33069)/looklook_usercenter?charset=utf8mb4&parseTime=true&loc=Local"
	t.Logf("使用数据库DSN: %s", dsn)

	// 手动设置Redis配置（Redis在Docker中，端口映射到36379）
	c.Cache = cache.CacheConf{
		cache.NodeConf{
			RedisConf: redis.RedisConf{
				Host: "localhost:36379",
				Pass: "G62m50oigInC30sf",
				Type: redis.NodeType,
			},
			Weight: 100,
		},
	}

	// 创建数据库连接
	conn := sqlx.NewMysql(dsn)

	testUserModel = model.NewUserModel(conn, c.Cache)
	testUserAuthModel = model.NewUserAuthModel(conn, c.Cache)

	t.Log("测试环境初始化完成")
}

// 清理测试数据
func cleanupTestData(t *testing.T, ctx context.Context, userIds []int64, authIds []int64) {
	for _, id := range userIds {
		_ = testUserModel.Delete(ctx, nil, id)
	}
	for _, id := range authIds {
		_ = testUserAuthModel.Delete(ctx, nil, id)
	}
	t.Log("测试数据清理完成")
}

// 测试1: Trans方法能否正常开启事务并提交
func TestTransCommit(t *testing.T) {
	setupTestDB(t)
	ctx := context.Background()

	var insertedUserId int64
	var insertedAuthId int64

	// 开启事务
	err := testUserModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 在事务中插入user
		user := &model.User{
			Mobile:     generateUniqueMobile(),
			Password:   "test_password",
			Nickname:   "测试用户_事务提交",
			Sex:        0,
			Avatar:     "http://example.com/avatar.jpg",
			Info:       "这是一个测试用户",
			DeleteTime: sql.NullTime{Time: time.Now(), Valid: true},
			DelState:   0,
			Version:    1,
		}

		result, err := testUserModel.Insert(ctx, session, user)
		if err != nil {
			return fmt.Errorf("插入user失败: %w", err)
		}

		insertedUserId, err = result.LastInsertId()
		if err != nil {
			return fmt.Errorf("获取user ID失败: %w", err)
		}

		t.Logf("在事务中插入user成功，ID: %d", insertedUserId)

		// 在事务中插入userAuth
		userAuth := &model.UserAuth{
			UserId:     insertedUserId,
			AuthKey:    fmt.Sprintf("test_key_%d", time.Now().Unix()),
			AuthType:   "test_type",
			DeleteTime: sql.NullTime{Time: time.Now(), Valid: true},
			DelState:   0,
			Version:    1,
		}

		result, err = testUserAuthModel.Insert(ctx, session, userAuth)
		if err != nil {
			return fmt.Errorf("插入userAuth失败: %w", err)
		}

		insertedAuthId, err = result.LastInsertId()
		if err != nil {
			return fmt.Errorf("获取userAuth ID失败: %w", err)
		}

		t.Logf("在事务中插入userAuth成功，ID: %d", insertedAuthId)

		// 正常返回，事务应该提交
		return nil
	})

	if err != nil {
		t.Fatalf("事务执行失败: %v", err)
	}

	t.Log("✓ 事务提交成功")

	// 验证数据是否真的插入了
	user, err := testUserModel.FindOne(ctx, insertedUserId)
	if err != nil {
		t.Fatalf("事务提交后查询user失败: %v", err)
	}
	if user.Id != insertedUserId {
		t.Fatalf("查询到的user ID不匹配，期望: %d, 实际: %d", insertedUserId, user.Id)
	}

	t.Log("✓ 事务提交后数据验证成功")

	// 清理测试数据
	defer cleanupTestData(t, ctx, []int64{insertedUserId}, []int64{insertedAuthId})
}

// 测试2: 事务回滚 - 遇到错误时应该回滚
func TestTransRollback(t *testing.T) {
	setupTestDB(t)
	ctx := context.Background()

	var attemptedUserId int64

	// 开启事务，故意在中间返回错误
	err := testUserModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 在事务中插入user
		user := &model.User{
			Mobile:     generateUniqueMobile(),
			Password:   "test_password",
			Nickname:   "测试用户_事务回滚",
			Sex:        0,
			Avatar:     "http://example.com/avatar.jpg",
			Info:       "这条数据应该被回滚",
			DeleteTime: sql.NullTime{Time: time.Now(), Valid: true},
			DelState:   0,
			Version:    1,
		}

		result, err := testUserModel.Insert(ctx, session, user)
		if err != nil {
			return fmt.Errorf("插入user失败: %w", err)
		}

		attemptedUserId, err = result.LastInsertId()
		if err != nil {
			return fmt.Errorf("获取user ID失败: %w", err)
		}

		t.Logf("在事务中插入user成功，ID: %d (应该被回滚)", attemptedUserId)

		// 故意返回错误，触发回滚
		return fmt.Errorf("模拟业务错误，触发事务回滚")
	})

	if err == nil {
		t.Fatal("期望事务返回错误，但实际成功了")
	}

	t.Logf("✓ 事务按预期返回错误: %v", err)

	// 验证数据是否被回滚了
	_, err = testUserModel.FindOne(ctx, attemptedUserId)
	if err != model.ErrNotFound {
		t.Fatalf("数据应该被回滚，但仍然能查询到，错误: %v", err)
	}

	t.Log("✓ 事务回滚成功，数据未插入")
}

// 测试3: 不使用事务 - session为nil
func TestWithoutTransaction(t *testing.T) {
	setupTestDB(t)
	ctx := context.Background()

	// 直接插入，不使用事务（session传nil）
	user := &model.User{
		Mobile:     generateUniqueMobile(),
		Password:   "test_password",
		Nickname:   "测试用户_无事务",
		Sex:        0,
		Avatar:     "http://example.com/avatar.jpg",
		Info:       "不使用事务插入的数据",
		DeleteTime: sql.NullTime{Time: time.Now(), Valid: true},
		DelState:   0,
		Version:    1,
	}

	result, err := testUserModel.Insert(ctx, nil, user)
	if err != nil {
		t.Fatalf("不使用事务插入失败: %v", err)
	}

	insertedUserId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("获取user ID失败: %v", err)
	}

	t.Logf("✓ 不使用事务插入成功，ID: %d", insertedUserId)

	// 验证数据
	foundUser, err := testUserModel.FindOne(ctx, insertedUserId)
	if err != nil {
		t.Fatalf("查询插入的数据失败: %v", err)
	}
	if foundUser.Nickname != user.Nickname {
		t.Fatalf("数据不匹配，期望昵称: %s, 实际: %s", user.Nickname, foundUser.Nickname)
	}

	t.Log("✓ 不使用事务的数据操作验证成功")

	// 清理测试数据
	defer cleanupTestData(t, ctx, []int64{insertedUserId}, []int64{})
}

// 测试4: 在事务中Update
func TestUpdateInTransaction(t *testing.T) {
	setupTestDB(t)
	ctx := context.Background()

	// 先插入一条数据
	user := &model.User{
		Mobile:     generateUniqueMobile(),
		Password:   "old_password",
		Nickname:   "旧昵称",
		Sex:        0,
		Avatar:     "http://example.com/old_avatar.jpg",
		Info:       "旧信息",
		DeleteTime: sql.NullTime{Time: time.Now(), Valid: true},
		DelState:   0,
		Version:    1,
	}

	result, err := testUserModel.Insert(ctx, nil, user)
	if err != nil {
		t.Fatalf("准备测试数据失败: %v", err)
	}

	insertedUserId, _ := result.LastInsertId()
	defer cleanupTestData(t, ctx, []int64{insertedUserId}, []int64{})

	// 在事务中更新
	err = testUserModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 先查询
		existingUser, err := testUserModel.FindOne(ctx, insertedUserId)
		if err != nil {
			return fmt.Errorf("查询失败: %w", err)
		}

		// 修改数据
		existingUser.Nickname = "新昵称_事务更新"
		existingUser.Password = "new_password"
		existingUser.Version = 2

		// 更新
		err = testUserModel.Update(ctx, session, existingUser)
		if err != nil {
			return fmt.Errorf("更新失败: %w", err)
		}

		t.Log("在事务中更新user成功")
		return nil
	})

	if err != nil {
		t.Fatalf("事务中更新失败: %v", err)
	}

	t.Log("✓ 事务更新提交成功")

	// 验证更新后的数据
	updatedUser, err := testUserModel.FindOne(ctx, insertedUserId)
	if err != nil {
		t.Fatalf("查询更新后的数据失败: %v", err)
	}

	if updatedUser.Nickname != "新昵称_事务更新" {
		t.Fatalf("昵称未更新，期望: %s, 实际: %s", "新昵称_事务更新", updatedUser.Nickname)
	}

	if updatedUser.Version != 2 {
		t.Fatalf("版本号未更新，期望: %d, 实际: %d", 2, updatedUser.Version)
	}

	t.Log("✓ 事务中Update功能验证成功")
}

// 测试5: 在事务中Delete
func TestDeleteInTransaction(t *testing.T) {
	setupTestDB(t)
	ctx := context.Background()

	// 先插入一条数据
	user := &model.User{
		Mobile:     generateUniqueMobile(),
		Password:   "test_password",
		Nickname:   "待删除用户",
		Sex:        0,
		Avatar:     "http://example.com/avatar.jpg",
		Info:       "这条数据将被删除",
		DeleteTime: sql.NullTime{Time: time.Now(), Valid: true},
		DelState:   0,
		Version:    1,
	}

	result, err := testUserModel.Insert(ctx, nil, user)
	if err != nil {
		t.Fatalf("准备测试数据失败: %v", err)
	}

	insertedUserId, _ := result.LastInsertId()

	// 在事务中删除
	err = testUserModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
		err := testUserModel.Delete(ctx, session, insertedUserId)
		if err != nil {
			return fmt.Errorf("删除失败: %w", err)
		}

		t.Log("在事务中删除user成功")
		return nil
	})

	if err != nil {
		t.Fatalf("事务中删除失败: %v", err)
	}

	t.Log("✓ 事务删除提交成功")

	// 验证数据已被删除
	_, err = testUserModel.FindOne(ctx, insertedUserId)
	if err != model.ErrNotFound {
		t.Fatalf("数据应该已被删除，但仍能查询到，错误: %v", err)
	}

	t.Log("✓ 事务中Delete功能验证成功")
}

// 测试6: 混合操作 - 在一个事务中进行Insert、Update、Delete
func TestMixedOperationsInTransaction(t *testing.T) {
	setupTestDB(t)
	ctx := context.Background()

	// 先准备一条待更新的数据
	existingUser := &model.User{
		Mobile:     generateUniqueMobile(),
		Password:   "old_password",
		Nickname:   "待更新用户",
		Sex:        0,
		Avatar:     "http://example.com/avatar.jpg",
		Info:       "将被更新",
		DeleteTime: sql.NullTime{Time: time.Now(), Valid: true},
		DelState:   0,
		Version:    1,
	}

	result, err := testUserModel.Insert(ctx, nil, existingUser)
	if err != nil {
		t.Fatalf("准备测试数据失败: %v", err)
	}
	existingUserId, _ := result.LastInsertId()

	// 再准备一条待删除的数据
	toDeleteUser := &model.User{
		Mobile:     generateUniqueMobile(),
		Password:   "test_password",
		Nickname:   "待删除用户",
		Sex:        0,
		Avatar:     "http://example.com/avatar.jpg",
		Info:       "将被删除",
		DeleteTime: sql.NullTime{Time: time.Now(), Valid: true},
		DelState:   0,
		Version:    1,
	}

	result, err = testUserModel.Insert(ctx, nil, toDeleteUser)
	if err != nil {
		t.Fatalf("准备测试数据失败: %v", err)
	}
	toDeleteUserId, _ := result.LastInsertId()

	var newUserId int64

	// 在一个事务中执行Insert、Update、Delete
	err = testUserModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 1. Insert 新用户
		newUser := &model.User{
			Mobile:     generateUniqueMobile(),
			Password:   "new_password",
			Nickname:   "事务中新建用户",
			Sex:        1,
			Avatar:     "http://example.com/new_avatar.jpg",
			Info:       "在事务中创建",
			DeleteTime: sql.NullTime{Time: time.Now(), Valid: true},
			DelState:   0,
			Version:    1,
		}

		result, err := testUserModel.Insert(ctx, session, newUser)
		if err != nil {
			return fmt.Errorf("插入新用户失败: %w", err)
		}
		newUserId, _ = result.LastInsertId()
		t.Logf("在事务中插入新用户，ID: %d", newUserId)

		// 2. Update 现有用户
		userToUpdate, err := testUserModel.FindOne(ctx, existingUserId)
		if err != nil {
			return fmt.Errorf("查询待更新用户失败: %w", err)
		}
		userToUpdate.Nickname = "已在事务中更新"
		userToUpdate.Version = 2

		err = testUserModel.Update(ctx, session, userToUpdate)
		if err != nil {
			return fmt.Errorf("更新用户失败: %w", err)
		}
		t.Logf("在事务中更新用户，ID: %d", existingUserId)

		// 3. Delete 用户
		err = testUserModel.Delete(ctx, session, toDeleteUserId)
		if err != nil {
			return fmt.Errorf("删除用户失败: %w", err)
		}
		t.Logf("在事务中删除用户，ID: %d", toDeleteUserId)

		return nil
	})

	if err != nil {
		t.Fatalf("混合操作事务失败: %v", err)
	}

	t.Log("✓ 混合操作事务提交成功")

	// 验证结果
	// 1. 验证新插入的用户存在
	newUser, err := testUserModel.FindOne(ctx, newUserId)
	if err != nil {
		t.Fatalf("查询新插入用户失败: %v", err)
	}
	if newUser.Nickname != "事务中新建用户" {
		t.Fatalf("新用户数据不正确")
	}
	t.Log("✓ 新插入的用户验证成功")

	// 2. 验证更新的用户
	updatedUser, err := testUserModel.FindOne(ctx, existingUserId)
	if err != nil {
		t.Fatalf("查询更新后用户失败: %v", err)
	}
	if updatedUser.Nickname != "已在事务中更新" || updatedUser.Version != 2 {
		t.Fatalf("用户更新不正确，昵称: %s, 版本: %d", updatedUser.Nickname, updatedUser.Version)
	}
	t.Log("✓ 更新的用户验证成功")

	// 3. 验证删除的用户
	_, err = testUserModel.FindOne(ctx, toDeleteUserId)
	if err != model.ErrNotFound {
		t.Fatalf("用户应该已被删除")
	}
	t.Log("✓ 删除的用户验证成功")

	// 清理测试数据
	defer cleanupTestData(t, ctx, []int64{newUserId, existingUserId}, []int64{})
}

// 测试7: 测试事务中的错误处理 - 部分操作失败应该全部回滚
func TestTransactionRollbackOnPartialFailure(t *testing.T) {
	setupTestDB(t)
	ctx := context.Background()

	var attemptedUserId int64
	var attemptedAuthId int64

	// 开启事务，插入user成功，插入userAuth成功后故意失败
	err := testUserModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 1. 插入user（应该成功）
		user := &model.User{
			Mobile:     generateUniqueMobile(),
			Password:   "test_password",
			Nickname:   "部分失败测试",
			Sex:        0,
			Avatar:     "http://example.com/avatar.jpg",
			Info:       "这条数据应该被回滚",
			DeleteTime: sql.NullTime{Time: time.Now(), Valid: true},
			DelState:   0,
			Version:    1,
		}

		result, err := testUserModel.Insert(ctx, session, user)
		if err != nil {
			return fmt.Errorf("插入user失败: %w", err)
		}
		attemptedUserId, _ = result.LastInsertId()
		t.Logf("在事务中插入user成功，ID: %d", attemptedUserId)

		// 2. 插入userAuth
		userAuth := &model.UserAuth{
			UserId:     attemptedUserId,
			AuthKey:    fmt.Sprintf("test_key_%d", time.Now().Unix()),
			AuthType:   "test_type",
			DeleteTime: sql.NullTime{Time: time.Now(), Valid: true},
			DelState:   0,
			Version:    1,
		}

		result, err = testUserAuthModel.Insert(ctx, session, userAuth)
		if err != nil {
			return fmt.Errorf("插入userAuth失败: %w", err)
		}
		attemptedAuthId, _ = result.LastInsertId()
		t.Logf("在事务中插入userAuth成功，ID: %d", attemptedAuthId)

		// 故意在这里返回错误，前面的两个操作应该都被回滚
		return fmt.Errorf("模拟业务逻辑错误")
	})

	if err == nil {
		t.Fatal("期望事务返回错误")
	}

	t.Logf("✓ 事务按预期失败: %v", err)

	// 验证user被回滚
	_, err = testUserModel.FindOne(ctx, attemptedUserId)
	if err != model.ErrNotFound {
		t.Fatalf("user应该被回滚，但仍能查询到")
	}

	// 验证userAuth被回滚
	_, err = testUserAuthModel.FindOne(ctx, attemptedAuthId)
	if err != model.ErrNotFound {
		t.Fatalf("userAuth应该被回滚，但仍能查询到")
	}

	t.Log("✓ 部分失败时的完全回滚验证成功")
}

// 测试8: 测试Update回滚
func TestUpdateRollback(t *testing.T) {
	setupTestDB(t)
	ctx := context.Background()

	// 先插入一条数据
	user := &model.User{
		Mobile:     generateUniqueMobile(),
		Password:   "original_password",
		Nickname:   "原始昵称",
		Sex:        0,
		Avatar:     "http://example.com/avatar.jpg",
		Info:       "原始信息",
		DeleteTime: sql.NullTime{Time: time.Now(), Valid: true},
		DelState:   0,
		Version:    1,
	}

	result, err := testUserModel.Insert(ctx, nil, user)
	if err != nil {
		t.Fatalf("准备测试数据失败: %v", err)
	}

	insertedUserId, _ := result.LastInsertId()
	defer cleanupTestData(t, ctx, []int64{insertedUserId}, []int64{})

	// 在事务中更新，然后回滚
	err = testUserModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 查询
		existingUser, err := testUserModel.FindOne(ctx, insertedUserId)
		if err != nil {
			return fmt.Errorf("查询失败: %w", err)
		}

		// 修改数据
		existingUser.Nickname = "这个修改应该被回滚"
		existingUser.Password = "rollback_password"
		existingUser.Version = 999

		// 更新
		err = testUserModel.Update(ctx, session, existingUser)
		if err != nil {
			return fmt.Errorf("更新失败: %w", err)
		}

		t.Log("在事务中更新user成功")

		// 故意返回错误，触发回滚
		return fmt.Errorf("触发回滚")
	})

	if err == nil {
		t.Fatal("期望事务返回错误")
	}

	t.Logf("✓ 事务按预期返回错误: %v", err)

	// 验证数据是否恢复原样
	unchangedUser, err := testUserModel.FindOne(ctx, insertedUserId)
	if err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}

	if unchangedUser.Nickname != "原始昵称" {
		t.Fatalf("Update应该被回滚，但昵称已改变。期望: %s, 实际: %s", "原始昵称", unchangedUser.Nickname)
	}

	if unchangedUser.Version != 1 {
		t.Fatalf("Update应该被回滚，但版本号已改变。期望: %d, 实际: %d", 1, unchangedUser.Version)
	}

	t.Log("✓ Update回滚验证成功")
}
