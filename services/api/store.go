package api

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Store interface {
	create(*gorm.DB, *models.Api) error
	update(*gorm.DB, *models.Api) (models.Api, error)
	deleteById(*gorm.DB, uint) error
	deleteByIds(*gorm.DB, []uint) error
	list(*gorm.DB, query.Query) (response.PageResult, error)
	getById(*gorm.DB, uint) (models.Api, error)
}

func ProvideStore() Store {
	return StoreImpl{}
}

type StoreImpl struct {
}

func (s StoreImpl) create(tx *gorm.DB, model *models.Api) error {
	return tx.Create(model).Error
}

func (s StoreImpl) update(tx *gorm.DB, partialUpdatedModel *models.Api) (updatedModel models.Api, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	result := tx.Updates(partialUpdatedModel)
	if result.Error != nil {
		return updatedModel, result.Error
	}
	if result.RowsAffected == 0 {
		return updatedModel, models.ErrModified("models.Api")
	}
	return s.getById(tx, partialUpdatedModel.Id)
}

func (s StoreImpl) deleteById(tx *gorm.DB, id uint) error {
	return tx.Delete(&models.Api{}, "id = ?", id).Error
}

func (s StoreImpl) deleteByIds(tx *gorm.DB, ids []uint) error {
	return tx.Delete(&models.Api{}, "id IN ?", ids).Error
}

func (s StoreImpl) list(tx *gorm.DB, apiQuery query.Query) (result response.PageResult, err error) {
	var total int64
	var list []models.Api
	if err = tx.Model(&models.Api{}).Count(&total).Error; err != nil {
		return
	}
	queryScope, err := query.MakeConditionFromQuery(apiQuery, models.Api{})
	if err != nil {
		return
	}
	if err = tx.
		Scopes(queryScope).
		Scopes(query.Paginate(apiQuery.Pagination)).
		Find(&list).Error; err != nil {
		return
	}
	result = response.PageResult{
		List:     list,
		Total:    total,
		PageNum:  apiQuery.Pagination.PageNum,
		PageSize: apiQuery.Pagination.PageSize,
	}
	return
}

func (s StoreImpl) getById(tx *gorm.DB, id uint) (model models.Api, err error) {
	err = tx.First(&model, "id = ?", id).Error
	return
}
