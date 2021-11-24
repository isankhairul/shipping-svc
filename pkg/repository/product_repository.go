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
	FindAll(ctx context.Context) (*[]entity.Product, error)
	Update(ctx context.Context, uid string, product *entity.Product) (*entity.Product, error)
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

func (repo *repo) FindAll(ctx context.Context) (*[]entity.Product, error) {
	var product []entity.Product
	err := repo.db.Find(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (repo *repo) Update(ctx context.Context, uid string, product *entity.Product) (*entity.Product, error) {
	err := repo.db.Where(&entity.Product{BaseIDModel: entity.BaseIDModel{UID: uid}}).Updates(product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}
