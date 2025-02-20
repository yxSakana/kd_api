package user

import (
	"context"

	"kd_api/api/user/v1"
	"kd_api/internal/logic/user"
)

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	tokenStr, err := user.Login(ctx, req.StudentId, req.Password)
	if err != nil {
		return nil, err
	}

	return &v1.LoginRes{Token: tokenStr}, nil
}
