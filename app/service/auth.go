package service

import (
	"context"

	"github.com/adamnasrudin03/go-market/app/dto"
	"github.com/adamnasrudin03/go-market/app/middlewares"
	"github.com/adamnasrudin03/go-market/app/models"
	"github.com/adamnasrudin03/go-market/app/repository"
	"github.com/adamnasrudin03/go-market/configs"
	"github.com/sirupsen/logrus"

	help "github.com/adamnasrudin03/go-helpers"
	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
)

type AuthService interface {
	Register(ctx context.Context, input dto.RegisterReq) (res *dto.UserRes, err error)
	Login(ctx context.Context, input dto.LoginReq) (res *dto.LoginRes, err error)
}

type AuthSrv struct {
	Repo      repository.AuthRepository
	RepoCache repository.CacheRepository
	Cfg       *configs.Configs
	Logger    *logrus.Logger
}

func NewAuthService(
	params AuthSrv,
) AuthService {
	return &params
}

func (srv *AuthSrv) Login(ctx context.Context, input dto.LoginReq) (res *dto.LoginRes, err error) {
	const opName = "AuthService-Login"
	defer help.PanicRecover(opName)

	user, err := srv.Repo.Login(ctx, input)
	if err != nil {
		srv.Logger.Errorf("%v error: %v", opName, err)
		return res, err
	}
	isExist := user != nil && user.ID > 0
	if !isExist {
		return nil, response_mapper.ErrDataNotFound("Pengguna", "User")
	}

	res = &dto.LoginRes{}
	res.Token, err = middlewares.GenerateToken(middlewares.JWTClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	})
	if err != nil {
		srv.Logger.Errorf("%v failed generate token: %v", opName, err)
		return res, err
	}

	return res, nil
}

func (srv *AuthSrv) Register(ctx context.Context, input dto.RegisterReq) (res *dto.UserRes, err error) {
	const opName = "AuthService-Register"
	defer help.PanicRecover(opName)

	user := models.User{
		Name:     input.Name,
		Password: input.Password,
		Email:    input.Email,
		Phone:    input.Phone,
	}

	err = srv.Repo.CheckIsDuplicate(ctx, dto.UserDetailReq{
		Email: input.Email,
		Phone: input.Phone,
	})
	if err != nil {
		return nil, err
	}

	detail, err := srv.Repo.Register(ctx, user)
	if err != nil || detail == nil {
		srv.Logger.Errorf("%v error create data: %v", opName, err)
		return nil, response_mapper.ErrCreatedDB()
	}

	return res.ConvertFromModel(*detail), nil
}
