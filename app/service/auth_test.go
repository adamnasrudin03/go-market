package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/adamnasrudin03/go-market/app/dto"
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
