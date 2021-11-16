package repository

import (
	"context"

	"gokit_example/pkg/entity"

	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type repo struct {
	db     *gorm.DB
	logger log.Logger
}

type ProductRepository interface {
	Save(ctx context.Context, product *entity.Product) (string, error)
	FindById(ctx context.Context, uid string) (*entity.Product, error)
}

func NewProductRepository(db *gorm.DB, logger log.Logger) ProductRepository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}

func (repo *repo) Save(ctx context.Context, product *entity.Product) (string, error) {
	err := repo.db.Create(&product).Error
	if err != nil {
		return "", err
	}

	return product.UID, nil

}

func (repo *repo) FindById(ctx context.Context, uid string) (*entity.Product, error) {
	var product entity.Product
	err := repo.db.Where(&entity.Product{BaseIDModel: entity.BaseIDModel{UID: uid}}).Find(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}
