package middle

import (
	"bytes"
	"io"
	"time"
	"yfapi/core/coreLog"
	"yfapi/core/coreSnowflake"
	error2 "yfapi/i18n/error"
	"yfapi/typedef"
	typedef_enum "yfapi/typedef/enum"
	"yfapi/typedef/response"
	"yfapi/util/easy"

	"github.com/gin-gonic/gin"
)

// RequestId
//
//	@Description:	获取请求id
//	@return			gin.HandlerFunc
func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := coreSnowflake.GetSnowId()
		c.Set("requestId", requestId)
		switch c.Request.Method {
		case "GET":
			// 获取查询参数并打印
			queryParams := c.Request.URL.Query()
			paramMap := make(map[string][]string)
			for key, values := range queryParams {
				paramMap[key] = values
			}
			coreLog.LogInfo("path:%s,request:%v", c.Request.URL.Path, easy.JSONStringFormObject(paramMap))
		default:
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				coreLog.LogError("RequestId ReadAll err:%+v", err)
				return
			}
			coreLog.LogInfo("path:%s,request:%s", c.Request.URL.Path, string(bodyBytes))
			// 恢复请求体，以便后续处理函数可以继续读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()
	}
}

func CheckHeaderData() gin.HandlerFunc {
	return func(context *gin.Context) {
		appVersion := context.GetHeader("appVersion")
		platform := context.GetHeader("platform")
		machineCode := context.GetHeader("machineCode")
		channel := context.GetHeader("channel")
		models := context.GetHeader("models")
		if len(appVersion) == 0 || len(platform) == 0 || len(channel) == 0 || len(machineCode) == 0 {
			coreLog.Error("CheckHeaderData err appVersion:%s,machineCode:%s,platform:%s,channel:%s", appVersion, machineCode, platform, channel)
			response.FailResponse(context, error2.I18nError{
				Code: error2.ErrorCodeParam,
				Msg:  nil,
			})
			context.Abort()
			return
		}
		if !easy.InArray(platform, typedef_enum.ClientTypeArray) {
			coreLog.Error("platform err platform:%s", platform)
			response.FailResponse(context, error2.I18nError{
				Code: error2.ErrorCodeParam,
				Msg:  nil,
			})
			context.Abort()
			return
		}
		var (
			address string
			Ip      = context.ClientIP()
		)
		if models == "" {
			models = platform
		}
		if !easy.IsPrivateIP(Ip) {
			address, _ = easy.GetIPLocation(Ip)
		}
		context.Set("headerData", typedef.HeaderData{
			AppVersion:  appVersion,
			Platform:    platform,
			Models:      models,
			MachineCode: machineCode,
			Channel:     channel,
			Address:     address,
			Ip:          Ip,
		})
	}
}

// CustomResponseWriter 封装 gin ResponseWriter 用于获取回包内容。
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := coreSnowflake.GetSnowId()
		c.Set("requestId", requestId)

		// 记录请求时间
		start := time.Now()

		// 使用自定义 ResponseWriter

		crw := &CustomResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = crw

		// 打印请求信息
		reqBody, _ := c.GetRawData()
		// 请求包体写回。
		if len(reqBody) > 0 {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}
		coreLog.LogInfo("【Request】:| %s | %s | %s | %s", c.ClientIP(), c.Request.Method, c.Request.RequestURI, reqBody)

		// 执行请求处理程序和其他中间件函数
		c.Next()

		// 记录回包内容和处理时间
		end := time.Now()
		latency := end.Sub(start)
		respBody := string(crw.body.Bytes())
		coreLog.LogInfo("【Response】:| %s | %s | %s | %s | (%v)\n", c.ClientIP(), c.Request.Method, c.Request.RequestURI, respBody, latency)
	}
}
