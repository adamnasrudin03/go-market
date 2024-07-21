package repository

import (
	"context"
	"errors"

	help "github.com/adamnasrudin03/go-helpers"
	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"github.com/adamnasrudin03/go-market/app/dto"
	"github.com/adamnasrudin03/go-market/app/models"
	"github.com/adamnasrudin03/go-market/configs"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Register(ctx context.Context, input models.User) (res *models.User, err error)
	Login(ctx context.Context, input dto.LoginReq) (res *models.User, er error)
	CheckIsDuplicate(ctx context.Context, input dto.UserDetailReq) (err error)
}

type AuthRepo struct {
	DB     *gorm.DB
	Cfg    *configs.Configs
	Logger *logrus.Logger
}

func NewAuthRepository(
	db *gorm.DB,
	cfg *configs.Configs,
	logger *logrus.Logger,
) AuthRepository {
	return &AuthRepo{
		DB:     db,
		Cfg:    cfg,
		Logger: logger,
	}
}
func (r *AuthRepo) Register(ctx context.Context, input models.User) (res *models.User, err error) {
	const opName = "AuthRepository-Register"
	err = r.DB.WithContext(ctx).Create(&input).Error
	if err != nil {
		r.Logger.Errorf("%v error register new user: %v ", opName, err)
		return nil, err
	}

	return &input, nil
}

func (r *AuthRepo) Login(ctx context.Context, input dto.LoginReq) (res *models.User, err error) {
	const opName = "AuthRepository-Login"
	err = r.DB.Select("id, phone, email, password").
		Where("email = ? OR phone = ?", input.Username, input.Username).
		WithContext(ctx).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		r.Logger.Errorf("%v error get db: %v ", opName, err)
		return nil, response_mapper.ErrDB()
	}

	if !help.PasswordIsValid(res.Password, input.Password) {
		r.Logger.Errorf("%v invalid password ", opName)
		return nil, response_mapper.ErrInvalid("Kata Sandi", "Password")
	}

	return res, nil
}

func (r *AuthRepo) CheckIsDuplicate(ctx context.Context, input dto.UserDetailReq) (err error) {
	if err = input.Validate(); err != nil {
		return err
	}

	var user *models.User
	req := dto.UserDetailReq{Columns: "id", NotID: input.NotID, Email: input.Email}
	if len(input.Email) > 0 {
		req.Email = input.Email
		user, _ = r.GetDetail(ctx, req)
		if user != nil && user.ID > 0 {
			return response_mapper.ErrIsDuplicate("Surel", "Email")
		}
		req.Email = ""
	}

	if len(input.Phone) > 0 {
		req.Phone = input.Phone
		user, _ = r.GetDetail(ctx, req)
		if user != nil && user.ID > 0 {
			return response_mapper.ErrIsDuplicate("Nomor Handphone", "Phone Number")
		}
		req.Phone = ""
	}

	return nil
}

func (r *AuthRepo) GetDetail(ctx context.Context, input dto.UserDetailReq) (res *models.User, err error) {
	const opName = "AuthRepository-GetDetail"

	column := "*"
	if input.Columns != "" {
		column = input.Columns
	}

	db := r.whereGetDetail(r.DB.Select(column), input)
	if err = db.WithContext(ctx).First(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		r.Logger.Errorf("%v error get db: %v ", opName, err)
		return nil, err
	}
	return res, nil
}

// whereGetDetail sets the where clause for the GetDetail function based on the
// input dto.
func (r *AuthRepo) whereGetDetail(db *gorm.DB, input dto.UserDetailReq) *gorm.DB {
	// Input is strongly typed
	if input.ID > 0 {
		db = db.Where("id = ?", input.ID)
	}
	if input.NotID > 0 {
		db = db.Where("id != ?", input.NotID)
	}
	if input.Email != "" {
		db = db.Where("email = ?", input.Email)
	}
	if input.Phone != "" {
		db = db.Where("Phone = ?", input.Phone)
	}
	if input.Name != "" {
		db = db.Where("name = ?", input.Name)
	}

	return db
}
