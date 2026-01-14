package internal

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wujunhui99/looklook/app/travel/model"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// TestConfig holds test configuration
type TestConfig struct {
	DBDataSource string
	RedisHost    string
	RedisPass    string
}

func getTestConfig() TestConfig {
	return TestConfig{
		DBDataSource: "root:PXDN93VRKUm8TeE7@tcp(127.0.0.1:33069)/looklook_travel?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai",
		RedisHost:    "127.0.0.1:36379",
		RedisPass:    "G62m50oigInC30sf",
	}
}

func getTestHomestayModel(t *testing.T) model.HomestayModel {
	cfg := getTestConfig()
	sqlConn := sqlx.NewMysql(cfg.DBDataSource)
	cacheConf := cache.CacheConf{
		{
			RedisConf: redis.RedisConf{
				Host: cfg.RedisHost,
				Pass: cfg.RedisPass,
				Type: "node",
			},
			Weight: 100,
		},
	}
	return model.NewHomestayModel(sqlConn, cacheConf)
}

// createTestHomestay creates a test homestay record
func createTestHomestay() *model.Homestay {
	return &model.Homestay{
		DeleteTime:          time.Now(),
		DelState:            model.DelStateNo,
		Version:             1,
		Title:               "Test Homestay",
		SubTitle:            "Test SubTitle",
		Banner:              "http://example.com/banner.jpg",
		Info:                "Test Info",
		PeopleNum:           4,
		HomestayBusinessId:  1,
		UserId:              1,
		RowState:            1,
		RowType:             0,
		FoodInfo:            "Test Food Info",
		FoodPrice:           1000,
		HomestayPrice:       20000,
		MarketHomestayPrice: 25000,
	}
}

// TestHomestayModel_InsertAndFindOne tests Insert and FindOne methods
func TestHomestayModel_InsertAndFindOne(t *testing.T) {
	homestayModel := getTestHomestayModel(t)
	ctx := context.Background()

	// Create test data
	homestay := createTestHomestay()
	homestay.Title = "TestInsertAndFindOne"

	// Insert
	result, err := homestayModel.Insert(ctx, nil, homestay)
	require.NoError(t, err, "Insert should not return error")

	id, err := result.LastInsertId()
	require.NoError(t, err, "LastInsertId should not return error")
	require.Greater(t, id, int64(0), "ID should be greater than 0")

	// FindOne
	found, err := homestayModel.FindOne(ctx, id)
	require.NoError(t, err, "FindOne should not return error")
	assert.Equal(t, homestay.Title, found.Title, "Title should match")
	assert.Equal(t, homestay.SubTitle, found.SubTitle, "SubTitle should match")
	assert.Equal(t, homestay.HomestayPrice, found.HomestayPrice, "HomestayPrice should match")

	// Cleanup: Hard delete the test data
	t.Cleanup(func() {
		_ = homestayModel.Delete(ctx, nil, id)
	})

	t.Logf("TestHomestayModel_InsertAndFindOne passed, id=%d", id)
}

// TestHomestayModel_Update tests Update method
func TestHomestayModel_Update(t *testing.T) {
	homestayModel := getTestHomestayModel(t)
	ctx := context.Background()

	// Insert test data
	homestay := createTestHomestay()
	homestay.Title = "TestUpdate_Before"
	result, err := homestayModel.Insert(ctx, nil, homestay)
	require.NoError(t, err)

	id, _ := result.LastInsertId()
	t.Cleanup(func() {
		_ = homestayModel.Delete(ctx, nil, id)
	})

	// Fetch the inserted record
	found, err := homestayModel.FindOne(ctx, id)
	require.NoError(t, err)

	// Update
	found.Title = "TestUpdate_After"
	found.HomestayPrice = 30000
	_, err = homestayModel.Update(ctx, nil, found)
	require.NoError(t, err, "Update should not return error")

	// Verify update
	updated, err := homestayModel.FindOne(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "TestUpdate_After", updated.Title, "Title should be updated")
	assert.Equal(t, int64(30000), updated.HomestayPrice, "HomestayPrice should be updated")

	t.Logf("TestHomestayModel_Update passed, id=%d", id)
}

// TestHomestayModel_UpdateWithVersion tests UpdateWithVersion method (optimistic locking)
func TestHomestayModel_UpdateWithVersion(t *testing.T) {
	homestayModel := getTestHomestayModel(t)
	ctx := context.Background()

	// Insert test data
	homestay := createTestHomestay()
	homestay.Title = "TestUpdateWithVersion"
	homestay.Version = 1
	result, err := homestayModel.Insert(ctx, nil, homestay)
	require.NoError(t, err)

	id, _ := result.LastInsertId()
	t.Cleanup(func() {
		_ = homestayModel.Delete(ctx, nil, id)
	})

	// Fetch the inserted record
	found, err := homestayModel.FindOne(ctx, id)
	require.NoError(t, err)
	originalVersion := found.Version

	// Update with version
	found.Title = "TestUpdateWithVersion_Updated"
	err = homestayModel.UpdateWithVersion(ctx, nil, found)
	require.NoError(t, err, "UpdateWithVersion should not return error")

	// Verify version was incremented
	updated, err := homestayModel.FindOne(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "TestUpdateWithVersion_Updated", updated.Title, "Title should be updated")
	assert.Equal(t, originalVersion+1, updated.Version, "Version should be incremented")

	// Test version conflict - try to update with old version
	found.Version = originalVersion // Use old version
	found.Title = "TestUpdateWithVersion_Conflict"
	err = homestayModel.UpdateWithVersion(ctx, nil, found)
	assert.Error(t, err, "UpdateWithVersion should return error for version conflict")
	assert.Equal(t, model.ErrNoRowsUpdate, err, "Error should be ErrNoRowsUpdate")

	t.Logf("TestHomestayModel_UpdateWithVersion passed, id=%d", id)
}

// TestHomestayModel_DeleteSoft tests DeleteSoft method
func TestHomestayModel_DeleteSoft(t *testing.T) {
	homestayModel := getTestHomestayModel(t)
	ctx := context.Background()

	// Insert test data
	homestay := createTestHomestay()
	homestay.Title = "TestDeleteSoft"
	result, err := homestayModel.Insert(ctx, nil, homestay)
	require.NoError(t, err)

	id, _ := result.LastInsertId()
	t.Cleanup(func() {
		_ = homestayModel.Delete(ctx, nil, id)
	})

	// Fetch the inserted record
	found, err := homestayModel.FindOne(ctx, id)
	require.NoError(t, err)

	// Soft delete
	err = homestayModel.DeleteSoft(ctx, nil, found)
	require.NoError(t, err, "DeleteSoft should not return error")

	// Verify soft delete
	deleted, err := homestayModel.FindOne(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, model.DelStateYes, deleted.DelState, "DelState should be DelStateYes")
	assert.True(t, deleted.DeleteTime.After(time.Now().Add(-time.Minute)), "DeleteTime should be recent")

	t.Logf("TestHomestayModel_DeleteSoft passed, id=%d", id)
}

// TestHomestayModel_FindSum tests FindSum method
func TestHomestayModel_FindSum(t *testing.T) {
	homestayModel := getTestHomestayModel(t)
	ctx := context.Background()

	// Insert multiple test data for sum calculation
	var ids []int64
	for i := 0; i < 3; i++ {
		homestay := createTestHomestay()
		homestay.Title = "TestFindSum"
		homestay.HomestayPrice = int64((i + 1) * 10000) // 10000, 20000, 30000
		result, err := homestayModel.Insert(ctx, nil, homestay)
		require.NoError(t, err)
		id, _ := result.LastInsertId()
		ids = append(ids, id)
	}

	t.Cleanup(func() {
		for _, id := range ids {
			_ = homestayModel.Delete(ctx, nil, id)
		}
	})

	// Test FindSum
	builder := homestayModel.SelectBuilder().Where("title = ?", "TestFindSum")
	sum, err := homestayModel.FindSum(ctx, builder, "homestay_price")
	require.NoError(t, err, "FindSum should not return error")
	assert.Equal(t, float64(60000), sum, "Sum should be 60000")

	// Test FindSum with empty field - should return error
	_, err = homestayModel.FindSum(ctx, builder, "")
	assert.Error(t, err, "FindSum with empty field should return error")

	t.Logf("TestHomestayModel_FindSum passed, sum=%v", sum)
}

// TestHomestayModel_FindCount tests FindCount method
func TestHomestayModel_FindCount(t *testing.T) {
	homestayModel := getTestHomestayModel(t)
	ctx := context.Background()

	// Insert multiple test data
	var ids []int64
	for i := 0; i < 5; i++ {
		homestay := createTestHomestay()
		homestay.Title = "TestFindCount"
		result, err := homestayModel.Insert(ctx, nil, homestay)
		require.NoError(t, err)
		id, _ := result.LastInsertId()
		ids = append(ids, id)
	}

	t.Cleanup(func() {
		for _, id := range ids {
			_ = homestayModel.Delete(ctx, nil, id)
		}
	})

	// Test FindCount
	builder := homestayModel.SelectBuilder().Where("title = ?", "TestFindCount")
	count, err := homestayModel.FindCount(ctx, builder, "id")
	require.NoError(t, err, "FindCount should not return error")
	assert.Equal(t, int64(5), count, "Count should be 5")

	// Test FindCount with empty field - should return error
	_, err = homestayModel.FindCount(ctx, builder, "")
	assert.Error(t, err, "FindCount with empty field should return error")

	t.Logf("TestHomestayModel_FindCount passed, count=%d", count)
}

// TestHomestayModel_FindAll tests FindAll method
func TestHomestayModel_FindAll(t *testing.T) {
	homestayModel := getTestHomestayModel(t)
	ctx := context.Background()

	// Insert multiple test data
	var ids []int64
	for i := 0; i < 3; i++ {
		homestay := createTestHomestay()
		homestay.Title = "TestFindAll"
		homestay.HomestayPrice = int64((i + 1) * 10000)
		result, err := homestayModel.Insert(ctx, nil, homestay)
		require.NoError(t, err)
		id, _ := result.LastInsertId()
		ids = append(ids, id)
	}

	t.Cleanup(func() {
		for _, id := range ids {
			_ = homestayModel.Delete(ctx, nil, id)
		}
	})

	// Test FindAll with default order (id DESC)
	builder := homestayModel.SelectBuilder().Where("title = ?", "TestFindAll")
	list, err := homestayModel.FindAll(ctx, builder, "")
	require.NoError(t, err, "FindAll should not return error")
	assert.Len(t, list, 3, "Should find 3 records")

	// Verify order is DESC by default
	for i := 0; i < len(list)-1; i++ {
		assert.Greater(t, list[i].Id, list[i+1].Id, "Records should be in DESC order")
	}

	// Test FindAll with custom order
	builder = homestayModel.SelectBuilder().Where("title = ?", "TestFindAll")
	list, err = homestayModel.FindAll(ctx, builder, "homestay_price ASC")
	require.NoError(t, err)
	assert.Len(t, list, 3)

	// Verify order is ASC by homestay_price
	for i := 0; i < len(list)-1; i++ {
		assert.Less(t, list[i].HomestayPrice, list[i+1].HomestayPrice, "Records should be in ASC order by price")
	}

	t.Logf("TestHomestayModel_FindAll passed, found %d records", len(list))
}

// TestHomestayModel_FindPageListByPage tests FindPageListByPage method
func TestHomestayModel_FindPageListByPage(t *testing.T) {
	homestayModel := getTestHomestayModel(t)
	ctx := context.Background()

	// Insert 10 test records
	var ids []int64
	for i := 0; i < 10; i++ {
		homestay := createTestHomestay()
		homestay.Title = "TestFindPageListByPage"
		result, err := homestayModel.Insert(ctx, nil, homestay)
		require.NoError(t, err)
		id, _ := result.LastInsertId()
		ids = append(ids, id)
	}

	t.Cleanup(func() {
		for _, id := range ids {
			_ = homestayModel.Delete(ctx, nil, id)
		}
	})

	// Test page 1
	builder := homestayModel.SelectBuilder().Where("title = ?", "TestFindPageListByPage")
	list, err := homestayModel.FindPageListByPage(ctx, builder, 1, 3, "")
	require.NoError(t, err, "FindPageListByPage should not return error")
	assert.Len(t, list, 3, "Page 1 should have 3 records")

	// Test page 2
	builder = homestayModel.SelectBuilder().Where("title = ?", "TestFindPageListByPage")
	list2, err := homestayModel.FindPageListByPage(ctx, builder, 2, 3, "")
	require.NoError(t, err)
	assert.Len(t, list2, 3, "Page 2 should have 3 records")

	// Ensure page 1 and page 2 have different records
	for _, item1 := range list {
		for _, item2 := range list2 {
			assert.NotEqual(t, item1.Id, item2.Id, "Page 1 and Page 2 should have different records")
		}
	}

	// Test page with page < 1 (should default to page 1)
	builder = homestayModel.SelectBuilder().Where("title = ?", "TestFindPageListByPage")
	listZero, err := homestayModel.FindPageListByPage(ctx, builder, 0, 3, "")
	require.NoError(t, err)
	assert.Len(t, listZero, 3, "Page 0 should default to page 1 and have 3 records")

	t.Logf("TestHomestayModel_FindPageListByPage passed")
}

// TestHomestayModel_FindPageListByPageWithTotal tests FindPageListByPageWithTotal method
func TestHomestayModel_FindPageListByPageWithTotal(t *testing.T) {
	homestayModel := getTestHomestayModel(t)
	ctx := context.Background()

	// Insert 7 test records
	var ids []int64
	for i := 0; i < 7; i++ {
		homestay := createTestHomestay()
		homestay.Title = "TestFindPageListByPageWithTotal"
		result, err := homestayModel.Insert(ctx, nil, homestay)
		require.NoError(t, err)
		id, _ := result.LastInsertId()
		ids = append(ids, id)
	}

	t.Cleanup(func() {
		for _, id := range ids {
			_ = homestayModel.Delete(ctx, nil, id)
		}
	})

	// Test with total
	builder := homestayModel.SelectBuilder().Where("title = ?", "TestFindPageListByPageWithTotal")
	list, total, err := homestayModel.FindPageListByPageWithTotal(ctx, builder, 1, 3, "")
	require.NoError(t, err, "FindPageListByPageWithTotal should not return error")
	assert.Len(t, list, 3, "Page 1 should have 3 records")
	assert.Equal(t, int64(7), total, "Total should be 7")

	// Test last page
	builder = homestayModel.SelectBuilder().Where("title = ?", "TestFindPageListByPageWithTotal")
	list, total, err = homestayModel.FindPageListByPageWithTotal(ctx, builder, 3, 3, "")
	require.NoError(t, err)
	assert.Len(t, list, 1, "Page 3 should have 1 record")
	assert.Equal(t, int64(7), total, "Total should still be 7")

	t.Logf("TestHomestayModel_FindPageListByPageWithTotal passed, total=%d", total)
}

// TestHomestayModel_FindPageListByIdDESC tests FindPageListByIdDESC method
func TestHomestayModel_FindPageListByIdDESC(t *testing.T) {
	homestayModel := getTestHomestayModel(t)
	ctx := context.Background()

	// Insert 5 test records
	var ids []int64
	for i := 0; i < 5; i++ {
		homestay := createTestHomestay()
		homestay.Title = "TestFindPageListByIdDESC"
		result, err := homestayModel.Insert(ctx, nil, homestay)
		require.NoError(t, err)
		id, _ := result.LastInsertId()
		ids = append(ids, id)
	}

	t.Cleanup(func() {
		for _, id := range ids {
			_ = homestayModel.Delete(ctx, nil, id)
		}
	})

	// First page (no preMinId)
	builder := homestayModel.SelectBuilder().Where("title = ?", "TestFindPageListByIdDESC")
	list, err := homestayModel.FindPageListByIdDESC(ctx, builder, 0, 2)
	require.NoError(t, err, "FindPageListByIdDESC should not return error")
	assert.Len(t, list, 2, "Should return 2 records")

	// Verify DESC order
	assert.Greater(t, list[0].Id, list[1].Id, "Records should be in DESC order")

	// Second page (with preMinId)
	minId := list[len(list)-1].Id
	builder = homestayModel.SelectBuilder().Where("title = ?", "TestFindPageListByIdDESC")
	list2, err := homestayModel.FindPageListByIdDESC(ctx, builder, minId, 2)
	require.NoError(t, err)
	assert.Len(t, list2, 2, "Should return 2 records")

	// All IDs in list2 should be less than minId
	for _, item := range list2 {
		assert.Less(t, item.Id, minId, "All IDs should be less than preMinId")
	}

	t.Logf("TestHomestayModel_FindPageListByIdDESC passed")
}

// TestHomestayModel_FindPageListByIdASC tests FindPageListByIdASC method
func TestHomestayModel_FindPageListByIdASC(t *testing.T) {
	homestayModel := getTestHomestayModel(t)
	ctx := context.Background()

	// Insert 5 test records
	var ids []int64
	for i := 0; i < 5; i++ {
		homestay := createTestHomestay()
		homestay.Title = "TestFindPageListByIdASC"
		result, err := homestayModel.Insert(ctx, nil, homestay)
		require.NoError(t, err)
		id, _ := result.LastInsertId()
		ids = append(ids, id)
	}

	t.Cleanup(func() {
		for _, id := range ids {
			_ = homestayModel.Delete(ctx, nil, id)
		}
	})

	// First page (no preMaxId)
	builder := homestayModel.SelectBuilder().Where("title = ?", "TestFindPageListByIdASC")
	list, err := homestayModel.FindPageListByIdASC(ctx, builder, 0, 2)
	require.NoError(t, err, "FindPageListByIdASC should not return error")
	assert.Len(t, list, 2, "Should return 2 records")

	// Verify ASC order
	assert.Less(t, list[0].Id, list[1].Id, "Records should be in ASC order")

	// Second page (with preMaxId)
	maxId := list[len(list)-1].Id
	builder = homestayModel.SelectBuilder().Where("title = ?", "TestFindPageListByIdASC")
	list2, err := homestayModel.FindPageListByIdASC(ctx, builder, maxId, 2)
	require.NoError(t, err)
	assert.Len(t, list2, 2, "Should return 2 records")

	// All IDs in list2 should be greater than maxId
	for _, item := range list2 {
		assert.Greater(t, item.Id, maxId, "All IDs should be greater than preMaxId")
	}

	t.Logf("TestHomestayModel_FindPageListByIdASC passed")
}

// TestHomestayModel_SelectBuilder tests SelectBuilder method
func TestHomestayModel_SelectBuilder(t *testing.T) {
	homestayModel := getTestHomestayModel(t)
	ctx := context.Background()

	// Insert test data
	homestay := createTestHomestay()
	homestay.Title = "TestSelectBuilder"
	homestay.HomestayPrice = 99999
	result, err := homestayModel.Insert(ctx, nil, homestay)
	require.NoError(t, err)

	id, _ := result.LastInsertId()
	t.Cleanup(func() {
		_ = homestayModel.Delete(ctx, nil, id)
	})

	// Use SelectBuilder to create custom query
	builder := homestayModel.SelectBuilder().
		Where("title = ?", "TestSelectBuilder").
		Where("homestay_price = ?", 99999)

	list, err := homestayModel.FindAll(ctx, builder, "")
	require.NoError(t, err, "Query with SelectBuilder should not return error")
	assert.Len(t, list, 1, "Should find 1 record")
	assert.Equal(t, "TestSelectBuilder", list[0].Title)
	assert.Equal(t, int64(99999), list[0].HomestayPrice)

	t.Logf("TestHomestayModel_SelectBuilder passed")
}

// TestHomestayModel_Trans tests Trans method (transaction)
func TestHomestayModel_Trans(t *testing.T) {
	homestayModel := getTestHomestayModel(t)
	ctx := context.Background()

	var insertedId int64

	// Test successful transaction
	err := homestayModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
		homestay := createTestHomestay()
		homestay.Title = "TestTrans"
		result, err := homestayModel.Insert(ctx, session, homestay)
		if err != nil {
			return err
		}
		insertedId, _ = result.LastInsertId()
		return nil
	})
	require.NoError(t, err, "Trans should not return error on success")

	// Verify data was inserted
	found, err := homestayModel.FindOne(ctx, insertedId)
	require.NoError(t, err)
	assert.Equal(t, "TestTrans", found.Title)

	// Cleanup
	_ = homestayModel.Delete(ctx, nil, insertedId)

	// Test rollback transaction
	var rollbackId int64
	err = homestayModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
		homestay := createTestHomestay()
		homestay.Title = "TestTrans_Rollback"
		result, err := homestayModel.Insert(ctx, session, homestay)
		if err != nil {
			return err
		}
		rollbackId, _ = result.LastInsertId()
		// Return error to trigger rollback
		return assert.AnError
	})
	assert.Error(t, err, "Trans should return error when function returns error")

	// Verify data was not inserted (rolled back)
	_, err = homestayModel.FindOne(ctx, rollbackId)
	assert.Equal(t, model.ErrNotFound, err, "Record should not exist after rollback")

	t.Logf("TestHomestayModel_Trans passed")
}
