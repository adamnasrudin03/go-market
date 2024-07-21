package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/adamnasrudin03/go-market/app/dto"
	"github.com/adamnasrudin03/go-market/app/models"
	"github.com/adamnasrudin03/go-market/app/repository/mocks"
	"github.com/adamnasrudin03/go-market/configs"
	"github.com/adamnasrudin03/go-market/pkg/driver"
	"github.com/stretchr/testify/suite"
)

type AuthServiceTestSuite struct {
	suite.Suite
	repo      *mocks.AuthRepository
	repoCache *mocks.CacheRepository
	ctx       context.Context
	service   AuthService
}

func (srv *AuthServiceTestSuite) SetupTest() {
	var (
		cfg    = configs.GetInstance()
		logger = driver.Logger(cfg)
	)
	srv.ctx = context.Background()

	srv.repo = &mocks.AuthRepository{}
	srv.repoCache = &mocks.CacheRepository{}
	params := AuthSrv{
		Repo:      srv.repo,
		RepoCache: srv.repoCache,
		Cfg:       cfg,
		Logger:    logger,
	}

	srv.service = NewAuthService(params)
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (srv *AuthServiceTestSuite) TestAuthSrv_Login() {
	tests := []struct {
		name     string
		params   dto.LoginReq
		mockFunc func(input dto.LoginReq)
		wantRes  *dto.LoginRes
		wantErr  bool
	}{
		{
			name:   "failed db",
			params: dto.LoginReq{},
			mockFunc: func(input dto.LoginReq) {
				srv.repo.On("Login", srv.ctx, input).Return(nil, errors.New("invalid")).Once()
			},
			wantRes: nil,
			wantErr: true,
		},
		{
			name:   "user not found",
			params: dto.LoginReq{},
			mockFunc: func(input dto.LoginReq) {
				srv.repo.On("Login", srv.ctx, input).Return(nil, nil).Once()
			},
			wantRes: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		srv.T().Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(tt.params)
			}
			gotRes, err := srv.service.Login(srv.ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthSrv.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("AuthSrv.Login() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func (srv *AuthServiceTestSuite) TestAuthSrv_Register() {
	params := dto.RegisterReq{
		Name:     "Adam",
		Email:    "adam@example.com",
		Phone:    "+6281234567890",
		Password: "secret123",
		Address:  "Kota Bekasi, Jawa barat, Indonesia",
	}
	tests := []struct {
		name     string
		req      dto.RegisterReq
		mockFunc func(input dto.RegisterReq)
		wantRes  *dto.UserRes
		wantErr  bool
	}{
		{
			name: "duplicate email",
			req:  params,
			mockFunc: func(input dto.RegisterReq) {
				srv.repo.On("CheckIsDuplicate", srv.ctx, dto.UserDetailReq{
					Email: input.Email,
					Phone: input.Phone,
				}).Return(errors.New("invalid email duplicate")).Once()
			},
			wantRes: nil,
			wantErr: true,
		},
		{
			name: "failed register user",
			req:  params,
			mockFunc: func(input dto.RegisterReq) {
				user := models.User{
					Name:     input.Name,
					Password: input.Password,
					Email:    input.Email,
					Phone:    input.Phone,
				}

				srv.repo.On("CheckIsDuplicate", srv.ctx, dto.UserDetailReq{
					Email: input.Email,
					Phone: input.Phone,
				}).Return(nil).Once()
				srv.repo.On("Register", srv.ctx, user).Return(nil, errors.New("invalid")).Once()
			},
			wantRes: nil,
			wantErr: true,
		},
		{
			name: "success",
			req:  params,
			mockFunc: func(input dto.RegisterReq) {
				user := models.User{
					Name:     input.Name,
					Password: input.Password,
					Email:    input.Email,
					Phone:    input.Phone,
				}

				srv.repo.On("CheckIsDuplicate", srv.ctx, dto.UserDetailReq{
					Email: input.Email,
					Phone: input.Phone,
				}).Return(nil).Once()
				srv.repo.On("Register", srv.ctx, user).Return(&models.User{
					ID:       101,
					Name:     input.Name,
					Password: input.Password,
					Email:    input.Email,
					Phone:    input.Phone,
				}, nil).Once()
			},
			wantRes: &dto.UserRes{
				ID:    101,
				Name:  params.Name,
				Phone: params.Phone,
				Email: params.Email,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		srv.T().Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(tt.req)
			}
			gotRes, err := srv.service.Register(srv.ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthSrv.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("AuthSrv.Register() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
