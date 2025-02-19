package WebCommonLogic

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type RespJson struct {
	Errcode int                    `json:"errcode"`
	Errmsg  string                 `json:"errmsg"`
	Data    map[string]interface{} `json:"data"`
}

func RenderDataJson(errcode int, errmsg string, Data map[string]interface{}, c *gin.Context) {
	var oRespJson RespJson
	oRespJson.Errcode = errcode
	oRespJson.Errmsg = errmsg
	oRespJson.Data = Data
	c.JSON(http.StatusOK, oRespJson)
	return
}

func RenderSuccessJson(c *gin.Context) {
	RenderDataJson(0, "", map[string]interface{}{}, c)
	return
}

func RenderErrorJson(errcode int, errmsg string, c *gin.Context) {
	var oRespJson RespJson
	oRespJson.Errcode = errcode
	oRespJson.Errmsg = errmsg
	c.JSON(http.StatusOK, oRespJson)
	return
}

// 结构体转义回参
func StructToMapViaReflect(obj interface{}) map[string]interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	result := make(map[string]interface{})
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		// 跳过未导出字段（PkgPath非空表示未导出）
		if field.PkgPath != "" {
			continue
		}

		// 获取字段值前先判断是否可导出
		val := v.Field(i)
		if val.CanInterface() {
			// 处理json tag
			tag := field.Tag.Get("json")
			if tag == "-" {
				continue // 明确忽略字段
			}
			name := strings.SplitN(tag, ",", 2)[0]
			if name == "" {
				name = field.Name
			}
			result[name] = val.Interface()
		}
	}
	return result
}

// Gin Request 参数解析
type WebGinContextParam struct {
	reqCtx *gin.Context
	init   bool
}

func NewGinCtxParam(c *gin.Context) *WebGinContextParam {
	return &WebGinContextParam{
		reqCtx: c,
		init:   true,
	}
}

func (wctx *WebGinContextParam) GetStringValue(name string, strVal string) string {
	if wctx.init != true {
		return strVal
	}

	tmpValue := wctx.reqCtx.DefaultQuery(name, strVal)

	return tmpValue
}

func (wctx *WebGinContextParam) GetIntValue(name string, nVal int32) int32 {
	var tmpValue string
	var ok bool
	if tmpValue, ok = wctx.reqCtx.GetQuery(name); !ok {
		return nVal
	}

	iRetVal, err := strconv.Atoi(tmpValue)
	if err != nil {
		return nVal
	}

	return int32(iRetVal)
}
