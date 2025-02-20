package cmd

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"kd_api/internal/controller/course"
	"kd_api/internal/controller/user"
	"kd_api/internal/logic/middleware"

	"kd_api/internal/spider/engine"
)

var (
	spiderCore = engine.NewKdClient()
	Main       = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Use(ReqLogHandler)
			s.Use(ErrorHandler)
			s.Group("/", func(group *ghttp.RouterGroup) {
				// Crawler Ctx
				group.Middleware(func(r *ghttp.Request) {
					r.SetCtxVar("spider", spiderCore)

					r.Middleware.Next()
				})
				// JSON Response Middleware
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				// Core
				group.Group("/v1", func(group *ghttp.RouterGroup) {
					// User
					group.Bind(
						user.NewV1(),
					)
					// Auth
					group.Group("/", func(group *ghttp.RouterGroup) {
						// Auth Middleware
						group.Middleware(middleware.Auth)
						// Course
						group.Bind(
							course.NewV1(),
						)
					})
				})
			})
			s.Run()
			return nil
		},
	}
)

var ReqLogHandler = func(r *ghttp.Request) {
	g.Log().Info(r.GetCtx(), r.URL)

	r.Middleware.Next()
}

var ErrorHandler = func(r *ghttp.Request) {
	r.Middleware.Next()

	//if r.Response.Status >= http.StatusInternalServerError {
	//	r.Response.ClearBuffer()
	//	r.Response.Writeln("Server internal error.")
	//}
	err := r.GetError()
	coverErr := recover()
	if err != nil || coverErr != nil {
		g.Log().Errorf(r.GetCtx(), "%v", err)
		r.Response.ClearBuffer()
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "Server internal error",
			Data:    nil,
		})
	}
}
