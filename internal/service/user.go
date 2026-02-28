package service

import (
	"context"
	"demoProject/dal/model"
	"demoProject/internal/api/dto"
	"demoProject/internal/e"
	"demoProject/internal/repo"
	"demoProject/pkg/util"
	"fmt"

	"go.uber.org/zap"
)
import "demoProject/pkg/snowflake"

// 实现IUserService接口
type UserService struct {
	userRepo repo.IUserRepo
}

// 接口类型本身就是一个“引用语义”的值，返回值不需要再取指针
func NewUserService(userRepo repo.IUserRepo) IUserService {
	return &UserService{userRepo: userRepo}

}

func (u *UserService) SignUp(ctx context.Context, req *dto.SignUpReq) error {
	// 1.判断用户存不存在
	is_exist, err := u.userRepo.UserIsExist(ctx, req.Username)
	if err != nil {
		return &util.CustomizedErr{
			Code: e.ERROR,
			Msg:  err.Error(),
		}
	}

	if is_exist {
		return &util.CustomizedErr{
			Code: e.ERROR_EXIST_USER,
			Msg:  e.GetMsg(e.ERROR_EXIST_USER),
		}
	}

	// 2.生成UUID
	uuid := snowflake.GenID()
	zap.L().Info(fmt.Sprintf("generate user ID:%s\n", uuid))

	// 3.保存用户信息至数据库
	u.userRepo.AddUser(ctx, &model.User{
		UserID:   uuid,
		Username: req.Username,
		Password: util.EncodeMD5(req.Password), // 敏感信息加密
		Email:    req.Email,
		Gender:   req.Gender,
	})

	return nil
}
