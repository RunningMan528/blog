package utils

import "gorm.io/gorm"

// 分页查询
func Paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 验证页码
		if page <= 0 {
			page = 1
		}
		// 验证size
		switch {
		case size > 100:
			size = 100
		case size < 10:
			size = 10
		}
		// 计算偏移
		offset := (page - 1) * size
		// Offset: 跳过N条记录
		// Limit: 返回最多 N 条记录
		return db.Offset(offset).Limit(size)
	}
}
