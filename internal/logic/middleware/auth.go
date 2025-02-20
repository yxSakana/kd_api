package middleware

import (
	"net/http"

	"github.com/gogf/gf/v2/net/ghttp"

	"kd_api/internal/logic/user"
)

func Auth(r *ghttp.Request) {
	tokenStr := r.GetHeader("Authorization")
	_, err := user.ParseToken(tokenStr)
	if err != nil {
		r.Response.WriteStatus(http.StatusUnauthorized)
		r.Exit()
	}

	r.Middleware.Next()
}
