package query

import (
	"fmt"

	"github.com/ronannnn/infra/utils"
	"gorm.io/gorm"
)

// Paginate for gorm pagination scopes
func Paginate(pageNum, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageSize <= 0 || pageNum <= 0 {
			return db
		}
		offset := pageSize * (pageNum - 1)
		return db.Offset(offset).Limit(pageSize)
	}
}

func MakeConditionFromQuery(query Query, setter QuerySetter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		condition := &DbConditionImpl{}
		ResolveQuery(query, setter, condition)
		return MakeCondition(condition)(db)
	}
}

func MakeCondition(condition *DbConditionImpl) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for k, v := range condition.Where {
			db = db.Where(k, v...)
		}
		for k, v := range condition.Or {
			db = db.Or(k, v...)
		}
		for k, v := range condition.Not {
			db = db.Not(k, v...)
		}
		for _, o := range condition.Order {
			db = db.Order(o)
		}
		if len(condition.Select) > 0 {
			db = db.Select(condition.Select)
		}
		return db
	}
}

func ResolveQueryRange(queryRange Range, fieldName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		start := queryRange.Start
		end := queryRange.End
		if !utils.IsZeroValue(start) {
			db.Where(fmt.Sprintf("%s >= ?", fieldName), start)
		}
		if !utils.IsZeroValue(end) {
			db.Where(fmt.Sprintf("%s <= ?", fieldName), end)
		}
		return db
	}
}
