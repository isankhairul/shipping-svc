package repository

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"strings"

	"gorm.io/gorm"
)

type productRepo struct {
	base BaseRepository
}

type ProductRepository interface {
	FindByUid(uid *string) (*entity.Product, error)
	FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.Product, *base.Pagination, error)
	Create(product *entity.Product) (*entity.Product, error)
	Update(uid string, input map[string]interface{}) error
	Delete(uid string) error
}

func NewProductRepository(br BaseRepository) ProductRepository {
	return &productRepo{br}
}

func (r *productRepo) FindByUid(uid *string) (*entity.Product, error) {
	var product entity.Product
	err := r.base.GetDB().
		Where("uid=?", uid).
		First(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepo) FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.Product, *base.Pagination, error) {
	var products []entity.Product
	var pagination base.Pagination

	query := r.base.GetDB()

	if filter["name"] != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(filter["name"].(string))+"%")
	}

	if filter["name"] != "" {
		query = query.Where("LOWER(sku) LIKE ?", "%"+strings.ToLower(filter["sku"].(string))+"%")
	}

	if filter["uom"] != "" {
		query = query.Where("LOWER(uom) = ?", strings.ToLower(filter["uom"].(string)))
	}

	if len(sort) > 0 {
		query = query.Order(sort)
	}

	pagination.Limit = limit
	pagination.Page = page
	err := query.Scopes(r.base.Paginate(products, &pagination, query, int64(len(products)))).
		Find(&products).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	return products, &pagination, nil
}

func (r *productRepo) Create(product *entity.Product) (*entity.Product, error) {
	err := r.base.GetDB().
		Create(product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepo) Update(uid string, input map[string]interface{}) error {
	err := r.base.GetDB().Model(&entity.Product{}).
		Where("uid=?", uid).
		Updates(input).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) Delete(uid string) error {
	var product entity.Product
	err := r.base.GetDB().
		Where("uid = ?", uid).
		Delete(&product).
		Error
	if err != nil {
		return err
	}

	return nil
}
