package service

import (
	"context"
	"time"

	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"github.com/adamnasrudin03/go-market/app/dto"
	"github.com/adamnasrudin03/go-market/app/models"
	"github.com/adamnasrudin03/go-market/app/repository"
	"github.com/adamnasrudin03/go-market/configs"
	"github.com/sirupsen/logrus"
)

type ProductService interface {
	GetByID(ctx context.Context, id uint64) (*models.Product, error)
	GetList(ctx context.Context, req dto.ProductListReq) (*response_mapper.Pagination, error)
}

type ProductSrv struct {
	Repo      repository.ProductRepository
	RepoCache repository.CacheRepository
	Cfg       *configs.Configs
	Logger    *logrus.Logger
}

func NewProductService(
	params ProductSrv,
) ProductService {
	return &params
}

func (s ProductSrv) GetByID(ctx context.Context, id uint64) (*models.Product, error) {
	var (
		opName = "ProductService-GetByID"
		err    error
		resp   models.Product
		key    = models.KeyCacheProductDetail(id)
	)

	ok := s.RepoCache.GetCache(ctx, key, &resp)
	if ok && resp.ID > 0 {
		return &resp, nil
	}

	detail, err := s.Repo.GetDetail(ctx, dto.ProductDetailReq{
		ID: id,
	})
	if err != nil {
		s.Logger.Errorf("%s, failed get detail: %v", opName, err)
		return nil, response_mapper.ErrDB()
	}

	isExist := detail != nil && detail.ID > 0
	if !isExist {
		return nil, response_mapper.ErrNotFound()
	}

	go s.RepoCache.CreateCache(context.Background(), key, detail, time.Minute)

	return detail, nil
}

func (s ProductSrv) GetList(ctx context.Context, req dto.ProductListReq) (*response_mapper.Pagination, error) {
	var (
		opName = "ProductService-GetList"
		err    error
		resp   *response_mapper.Pagination
	)
	err = req.Validate()
	if err != nil {
		return nil, err
	}

	data, err := s.Repo.GetList(ctx, req)
	if err != nil {
		s.Logger.Errorf("%s, failed get list: %v", opName, err)
		return nil, response_mapper.ErrDB()
	}

	totalRecords := len(data)
	resp = &response_mapper.Pagination{
		Data: data,
		Meta: response_mapper.Meta{
			Page:         req.Page,
			Limit:        req.Limit,
			TotalRecords: totalRecords,
		},
	}

	// total records in less than limit
	if totalRecords > 0 && totalRecords != req.Limit {
		return resp, nil
	}

	// get total data
	if totalRecords > 0 {
		req.CustomColumns = "id"
		req.IsNotDefaultQuery = true
		req.Offset = (req.Page - 1) * req.Limit
		req.Limit = models.DefaultLimitIsTotalDataTrue * req.Limit

		total, err := s.Repo.GetList(ctx, req)
		if err != nil {
			s.Logger.Errorf("%s, failed get total data: %v", opName, err)
			return nil, response_mapper.ErrDB()
		}
		totalRecords = len(total)
		resp.Meta.TotalRecords = totalRecords
	}

	return resp, nil
}
