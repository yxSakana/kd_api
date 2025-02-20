package v1

import "github.com/gogf/gf/v2/frame/g"

type LoginReq struct {
	g.Meta    `path:"user/login" method:"post" tags:"user"`
	StudentId string `json:"student_id" v:"required"`
	Password  string `json:"password" v:"required"`
}
type LoginRes struct {
	Token string `json:"token"`
}
