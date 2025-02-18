package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qwe826344858/mygin/Route"
)

func main() {
	r := gin.Default()
	Route.Register(r)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
