package main

import (
	_ "kd_api/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gctx"

	"kd_api/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
