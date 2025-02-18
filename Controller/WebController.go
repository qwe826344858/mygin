package Controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	extProto "github.com/qwe826344858/dockerGoProject/ExternalProto"
	proto "github.com/qwe826344858/dockerGoProject/proto"
	wclg "github.com/qwe826344858/mygin/WebCommonLogic"
)

func Ping(c *gin.Context) {
	Data := map[string]interface{}{
		"message": "pong",
	}
	wclg.RenderDataJson(0, "", Data, c)
	return
}

// 调用go服务查询
func GetSteamItemInfoByGoAo(c *gin.Context) {
	f, _, err := proto.GetDockerGoProjectAoClient()
	defer f.CloseClient()
	if err != nil {
		fmt.Printf("GetSteamItemInfoByGoAo::GetDockerGoProjectAoClient err:%v", err)
		wclg.RenderErrorJson(20001, "系统繁忙", c)
		return
	}

	//req := &proto.GetItemInfoReq{
	//	ReqHeader: &proto.RequestHeader{},
	//	ItemId:    1,
	//}
	//
	//resp, err := client.GetItemInfo(context.TODO(), req)
	//if err != nil {
	//	wclg.RenderErrorJson(20002, "服务异常", c)
	//	return
	//}

	resp := new(proto.GetItemInfoResp)
	//返回结果
	wclg.RenderDataJson(0, "", wclg.StructToMapViaReflect(resp), c)
	return
}

// 调用python服务查询
func GetSteamItemInfoByPythonAo(c *gin.Context) {
	f, dockerPyClient, err := extProto.GetDockerProjectAoClient()
	defer f.CloseClient()
	if err != nil {
		fmt.Printf("GetSteamItemInfoByPythonAo::GetDockerProjectAoClient err:%v", err)
		wclg.RenderErrorJson(20001, "系统繁忙", c)
		return
	}

	req := &extProto.GetItemInfoReq{
		ReqHeader: &extProto.RequestHeader{},
		ItemId:    2,
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
