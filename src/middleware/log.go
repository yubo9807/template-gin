package middleware

import (
	"fmt"
	"io"
	"runtime/debug"
	"server/configs"
	"server/src/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Log(c *gin.Context) {
	defer func() {
		if msg := recover(); msg != nil {
			debug.PrintStack()
			msgStr := fmt.Sprintf("%v", msg)
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  msgStr,
			})
			c.Abort()
			logWrite(c, "panic: "+msgStr+"\n"+string(debug.Stack()))
		}
	}()

	c.Next()

	logWrite(c, "result: "+ContextGet(c).Result)
}

func logWrite(c *gin.Context, msg string) {
	ctx := ContextGet(c)
	t := strings.ReplaceAll(ctx.RequestTime.String()[0:26], "-", "/")

	ip := c.ClientIP()
	method := c.Request.Method
	path := c.Request.RequestURI

	header := "\nhearers: "
	for k, v := range c.Request.Header {
		header += "\n\t" + k + ":" + v[0]
	}

	bodyData, _ := io.ReadAll(c.Request.Body)
	body := string(bodyData)

	content := t + "\t" + ip + "\t" + method + "\t" + path
	runTime := time.Since(ctx.RequestTime).String()
	content += "\nstatus:" + strconv.Itoa(c.Writer.Status()) + "\trun_time:" + runTime
	// content += header
	if body != "" {
		content += "\nbody: " + body
	}
	if msg != "" {
		content += "\n" + msg
	}

	utils.LogWrite(configs.Config.LogDir, content)
	utils.LogClear(configs.Config.LogDir, configs.Config.LogReserveTime)
}
