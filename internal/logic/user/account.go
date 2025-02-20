package user

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"

	"kd_api/internal/dao"
	"kd_api/internal/model/do"
	"kd_api/internal/model/entity"
	"kd_api/internal/spider/engine"
)

func Login(ctx context.Context, studentId, password string,
) (tokenStr string, err error) {
	// Get dean office system session id
	// 1) From Redis cache

	// 2) Login dean office system again
	crawler, ok := ctx.Value("spider").(*engine.KdClient)
	if !ok {
		return "", gerror.New("Error from context to get spider")
	}
	err = crawler.LoginSystem(studentId, password, 3, false)
	if err != nil {
		return "", err
	}
	studentInfo, err := crawler.GetStudentInfo()
	if err != nil {
		return "", err
	}

	// Create and return token by studentId
	tokenStr, err = GenerateToken(studentInfo.StudentId)
	if err != nil {
		return "", err
	}

	_, err = dao.User.Ctx(ctx).Data(do.User{
		StudentId:   studentInfo.StudentId,
		Password:    gbase64.EncodeString(password),
		StudentCode: studentInfo.StudentCode,
		Name:        studentInfo.Name,
		SchoolCode:  studentInfo.SchoolCode,
		SchoolName:  studentInfo.SchoolName,
	}).Save()
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func GetStudentID(ctx context.Context) (uint32, error) {
	authHeader := g.RequestFromCtx(ctx).GetHeader("Authorization")
	claims, err := ParseToken(authHeader)
	if err != nil {
		return 0, err
	}
	return claims.StudentId, nil
}

func GetStudentInfo(ctx context.Context) (entity.User, error) {
	studentId, err := GetStudentID(ctx)
	if err != nil {
		return entity.User{}, err
	}

	studentIdInt := gconv.Uint(studentId)
	var userInfo entity.User
	r, err := dao.User.Ctx(ctx).WherePri(studentIdInt).One()
	if err != nil {
		return entity.User{}, err
	}
	err = r.Struct(&userInfo)
	if err != nil {
		return entity.User{}, err
	}
	return userInfo, nil
}
