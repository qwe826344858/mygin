package Controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	extProto "github.com/qwe826344858/dockerGoProject/ExternalProto"
	proto "github.com/qwe826344858/dockerGoProject/proto"
	wclg "github.com/qwe826344858/mygin/WebCommonLogic"
	"time"
)

func Ping(c *gin.Context) {
	Data := map[string]interface{}{
		"message": "pong",
	}
	fmt.Printf("有人Ping! now:%v \n", time.Now())
	wclg.RenderDataJson(0, "", Data, c)
	return
}

// 调用go服务查询
func GetSteamItemInfoByGoAo(c *gin.Context) {
	//解析入参:
	param := wclg.NewGinCtxParam(c)
	fmt.Printf("param:%v", param)
	litemId := param.GetIntValue("item_id", 0)
	if litemId < 1 {
		wclg.RenderErrorJson(30001, "请输入商品id", c)
		return
	}

	f, client, err := proto.GetDockerGoProjectAoClient()
	defer f.CloseClient()
	if err != nil {
		fmt.Printf("GetSteamItemInfoByGoAo::GetDockerGoProjectAoClient err:%v", err)
		wclg.RenderErrorJson(20001, "系统繁忙", c)
		return
	}

	req := &proto.GetItemInfoReq{
		ReqHeader: &proto.RequestHeader{},
		ItemId:    int64(litemId),
	}

	resp, err := client.GetItemInfo(context.TODO(), req)
	if err != nil {
		wclg.RenderErrorJson(20002, "服务异常", c)
		return
	}

	//返回结果
	wclg.RenderDataJson(0, "", wclg.StructToMapViaReflect(resp), c)
	return
}

// 调用python服务查询
func GetSteamItemInfoByPythonAo(c *gin.Context) {
	//解析入参:
	param := wclg.NewGinCtxParam(c)
	fmt.Printf("param:%v", param)
	litemId := param.GetIntValue("item_id", 0)
	if litemId < 1 {
		wclg.RenderErrorJson(30001, "请输入商品id", c)
		return
	}

	f, dockerPyClient, err := extProto.GetDockerProjectAoClient()
	defer f.CloseClient()
	if err != nil {
		fmt.Printf("GetSteamItemInfoByPythonAo::GetDockerProjectAoClient err:%v", err)
		wclg.RenderErrorJson(20001, "系统繁忙", c)
		return
	}

	req := &extProto.GetItemInfoReq{
		ReqHeader: &extProto.RequestHeader{},
		ItemId:    int64(litemId),
	}
	resp, err := dockerPyClient.GetItemInfo(context.TODO(), req)
	if err != nil {
		wclg.RenderErrorJson(20002, "服务异常", c)
		return
	}

	// 返回结果
	wclg.RenderDataJson(0, "", wclg.StructToMapViaReflect(resp), c)
	return
}

func PanicTest(c *gin.Context) {
	var v interface{}
	fmt.Printf("我要panic了! now:%v \n", time.Now())
	v = "panic test"
	panic(v)
	return
}
