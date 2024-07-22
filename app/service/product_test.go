package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"github.com/adamnasrudin03/go-market/app/dto"
	"github.com/adamnasrudin03/go-market/app/models"
	"github.com/adamnasrudin03/go-market/app/repository/mocks"
	"github.com/adamnasrudin03/go-market/configs"
	"github.com/adamnasrudin03/go-market/pkg/driver"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProductServiceTestSuite struct {
	suite.Suite
	repo      *mocks.ProductRepository
	repoCache *mocks.CacheRepository
	ctx       context.Context
	service   ProductService
	product   models.Product
	products  []models.Product
}

func (srv *ProductServiceTestSuite) SetupTest() {
	var (
		cfg    = configs.GetInstance()
		logger = driver.Logger(cfg)
	)
	srv.product = models.Product{
		ID:          1,
		Name:        "Product A",
		WarehouseID: 2,
	}
	srv.products = []models.Product{
		srv.product,
		{
			ID:          2,
			Name:        "Product B",
			WarehouseID: 1,
		},
	}

	srv.repo = &mocks.ProductRepository{}
	srv.repoCache = &mocks.CacheRepository{}
	srv.ctx = context.Background()
	srv.service = NewProductService(ProductSrv{
		Repo:      srv.repo,
		RepoCache: srv.repoCache,
		Cfg:       cfg,
		Logger:    logger,
	})
}

func TestProductService(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}

func (srv *ProductServiceTestSuite) TestProductSrv_GetByID() {
	tests := []struct {
		name     string
		id       uint64
		mockFunc func(input uint64)
		want     *models.Product
		wantErr  bool
	}{
		{
			name: "Success with cache",
			id:   srv.product.ID,
			mockFunc: func(input uint64) {
				key := models.KeyCacheProductDetail(input)
				res := srv.product
				srv.repoCache.On("GetCache", mock.Anything, key, &models.Product{
					ID: 0,
				}).Return(true).Run(func(args mock.Arguments) {
					target := args.Get(2).(*models.Product)
					*target = res
				}).Once()
			},
			want:    &srv.product,
			wantErr: false,
		},
		{
			name: "failed ge db",
			id:   srv.product.ID,
			mockFunc: func(input uint64) {
				key := models.KeyCacheProductDetail(input)
				srv.repoCache.On("GetCache", mock.Anything, key, &models.Product{
					ID: 0,
				}).Return(false).Once()

				srv.repo.On("GetDetail", mock.Anything, dto.ProductDetailReq{
					ID: input,
				}).Return(nil, errors.New("invalid")).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not found",
			id:   srv.product.ID,
			mockFunc: func(input uint64) {
				key := models.KeyCacheProductDetail(input)
				srv.repoCache.On("GetCache", mock.Anything, key, &models.Product{
					ID: 0,
				}).Return(false).Once()

				srv.repo.On("GetDetail", mock.Anything, dto.ProductDetailReq{
					ID: input,
				}).Return(nil, nil).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			id:   srv.product.ID,
			mockFunc: func(input uint64) {
				key := models.KeyCacheProductDetail(input)
				srv.repoCache.On("GetCache", mock.Anything, key, &models.Product{
					ID: 0,
				}).Return(false).Once()

				srv.repo.On("GetDetail", mock.Anything, dto.ProductDetailReq{
					ID: input,
				}).Return(&srv.product, nil).Once()

				srv.repoCache.On("CreateCache", mock.Anything, key, &srv.product, time.Minute).Return().Once()

			},
			want:    &srv.product,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		srv.T().Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(tt.id)
			}

			got, err := srv.service.GetByID(srv.ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductSrv.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductSrv.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (srv *ProductServiceTestSuite) TestProductSrv_GetList() {
	params := dto.ProductListReq{
		BasedFilter: models.BasedFilter{
			Limit: 1,
			Page:  1,
		},
	}

	tests := []struct {
		name     string
		req      dto.ProductListReq
		mockFunc func(input dto.ProductListReq)
		want     *response_mapper.Pagination
		wantErr  bool
	}{
		{
			name: "invalid params",
			req: dto.ProductListReq{
				BasedFilter: models.BasedFilter{
					OrderBy: "invalid",
				},
			},
			mockFunc: func(input dto.ProductListReq) {
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed get records",
			req:  params,
			mockFunc: func(input dto.ProductListReq) {
				srv.repo.On("GetList", mock.Anything, input).Return(nil, errors.New("invalid")).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success total records less than limit",
			req: dto.ProductListReq{
				BasedFilter: models.BasedFilter{
					Limit: 10,
					Page:  1,
				},
			},
			mockFunc: func(input dto.ProductListReq) {
				srv.repo.On("GetList", mock.Anything, input).Return(srv.products, nil).Once()
			},
			want: &response_mapper.Pagination{
				Meta: response_mapper.Meta{
					Page:         1,
					Limit:        10,
					TotalRecords: len(srv.products),
				},
				Data: srv.products,
			},
			wantErr: false,
		},
		{
			name: "failed get total records",
			req:  params,
			mockFunc: func(input dto.ProductListReq) {
				srv.repo.On("GetList", mock.Anything, input).Return([]models.Product{srv.products[0]}, nil).Once()

				input.CustomColumns = "id"
				input.IsNotDefaultQuery = true
				input.Offset = (input.Page - 1) * input.Limit
				input.Limit = models.DefaultLimitIsTotalDataTrue * input.Limit
				srv.repo.On("GetList", mock.Anything, input).Return(nil, errors.New("invalid")).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success total records more than limit",
			req:  params,
			mockFunc: func(input dto.ProductListReq) {
				srv.repo.On("GetList", mock.Anything, input).Return([]models.Product{srv.products[0]}, nil).Once()

				input.CustomColumns = "id"
				input.IsNotDefaultQuery = true
				input.Offset = (input.Page - 1) * input.Limit
				input.Limit = models.DefaultLimitIsTotalDataTrue * input.Limit

				total := []models.Product{}
				for i := 0; i < len(srv.products); i++ {
					val := models.Product{ID: srv.products[i].ID}
					total = append(total, val)
				}
				srv.repo.On("GetList", mock.Anything, input).Return(total, nil).Once()
			},
			want: &response_mapper.Pagination{
				Meta: response_mapper.Meta{
					Page:         params.Page,
					Limit:        params.Limit,
					TotalRecords: len(srv.products),
				},
				Data: []models.Product{srv.products[0]},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		srv.T().Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(tt.req)
			}

			got, err := srv.service.GetList(srv.ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductSrv.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductSrv.GetList() = %v, want %v", got, tt.want)
			}
		})
	}
}
