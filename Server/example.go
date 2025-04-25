package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qwe826344858/mygin/Route"
	"log"
	"net/http"
)

func main() {
	// 默认使用 Recovery 和 Logger 中间件
	r := gin.Default()
	//r.Use(gin.Recovery()) // 中间件panic捕获
	//r.Use(gin.Logger()) //  gin log 日志相关

	// r.Use(MiddleWare()) // 自定义中间件

	Route.Register(r)
	if err := r.Run(); err != nil { // 监听并在 0.0.0.0:8080 上启动服务
		log.Fatalf("Failed to start server: %v", err)
	}
}

// 自定义中间件
func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("中间件启动\n")
		sessionId, _ := c.Cookie("sessionId")
		if sessionId == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "sessionId 为空或未找到",
			})
		}
		c.Next()
		fmt.Printf("中间件结束\n")
	}
}
