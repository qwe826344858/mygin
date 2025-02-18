package Route

import (
	"github.com/gin-gonic/gin"
	"github.com/qwe826344858/mygin/Controller"
)

func Register(r *gin.Engine) {
	r.GET("/ping", Controller.Ping)
	r.GET("/getsteamiteminfogo", Controller.GetSteamItemInfoByGoAo)
	r.GET("/getsteamiteminfopy", Controller.GetSteamItemInfoByPythonAo)

}
