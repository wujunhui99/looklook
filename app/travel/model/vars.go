package model

import (
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound
var ErrNoRowsUpdate = errors.New("update db no rows change")

// Soft delete state constants
var (
	DelStateNo  int64 = 0 // 未删除
	DelStateYes int64 = 1 // 已删除
)

// HomestayActivity row_type constants
const (
	HomestayActivityPreferredType = "preferred" // Preferred homestays
	HomestayActivityGoodBusiType  = "goodBoss"  // Good business owners
)

// HomestayActivity row_status constants
const (
	HomestayActivityUpStatus   = "up"   // Active/Published
	HomestayActivityDownStatus = "down" // Inactive/Unpublished
)
