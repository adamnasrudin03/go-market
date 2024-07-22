package repository

import (
	"context"
	"errors"

	"github.com/adamnasrudin03/go-market/app/dto"
	"github.com/adamnasrudin03/go-market/app/models"
	"github.com/adamnasrudin03/go-market/configs"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetDetail(ctx context.Context, req dto.ProductDetailReq) (*models.Product, error)
	GetList(ctx context.Context, req dto.ProductListReq) ([]models.Product, error)
}

type ProductRepo struct {
	DB     *gorm.DB
	Cfg    *configs.Configs
	Logger *logrus.Logger
}

func NewProductRepository(
	db *gorm.DB,
	cfg *configs.Configs,
	logger *logrus.Logger,
) ProductRepository {
	return &ProductRepo{
		DB:     db,
		Cfg:    cfg,
		Logger: logger,
	}
}

func (r *ProductRepo) GetDetail(ctx context.Context, req dto.ProductDetailReq) (*models.Product, error) {
	var (
		opName = "ProductRepository-GetDetail"
		err    error
		resp   *models.Product
		column = "*"
	)
	if req.CustomColumn != "" {
		column = req.CustomColumn
	}

	db := r.DB.WithContext(ctx).Model(&models.Product{}).Select(column)
	if req.ID > 0 {
		db = db.Where("id = ?", req.ID)
	}
	if req.WarehouseID > 0 {
		db = db.Where("warehouse_id = ?", req.WarehouseID)
	}

	err = db.First(&resp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		r.Logger.Errorf("%v error: %v ", opName, err)
		return nil, err
	}

	return resp, nil
}

func (r *ProductRepo) GetList(ctx context.Context, req dto.ProductListReq) ([]models.Product, error) {
	var (
		opName = "ProductRepository-GetList"
		err    error
		resp   []models.Product
		column = "*"
	)
	if req.CustomColumns != "" {
		column = req.CustomColumns
	}

	db := r.DB.WithContext(ctx).Model(&models.Product{}).Select(column)
	if req.Search != "" {
		db = db.Where("name LIKE ?", "%"+req.Search+"%")
	}

	if !req.IsNotDefaultQuery {
		req.BasedFilter = req.DefaultQuery()
	}
	if !req.IsNoLimit {
		db = db.Offset(int(req.Offset)).Limit(int(req.Limit))
	}

	if models.IsValidOrderBy[req.OrderBy] && req.SortBy != "" {
		db = db.Order(req.SortBy + " " + req.OrderBy)
	}

	err = db.Find(&resp).Error
	if err != nil {
		r.Logger.Errorf("%v error: %v ", opName, err)
		return nil, err
	}

	return resp, nil
}
