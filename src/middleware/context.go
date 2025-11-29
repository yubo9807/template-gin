package middleware

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

type ContextCustom struct {
	*gin.Context
	RequestTime time.Time
	Result      string
}

func ContextMiddleware(c *gin.Context) {
	customCtx := &ContextCustom{Context: c}
	customCtx.RequestTime = time.Now()
	c.Set("customCtx", customCtx)
	c.Next()
}

func ContextGet(c *gin.Context) *ContextCustom {
	value, exist := c.Get("customCtx")
	if !exist {
		return nil
	}
	ctx, ok := value.(*ContextCustom)
	if !ok {
		return nil
	}
	return ctx
}

func (c *ContextCustom) Success() {
	c.SuccessData("")
}

func (c *ContextCustom) SuccessData(data interface{}) {
	m := gin.H{
		"code":    200,
		"data":    data,
		"msg":     "success",
		"runTime": time.Since(c.RequestTime).String(),
	}
	c.JSON(200, m)
	c.Abort()

	jsonBytes, _ := json.Marshal(m)
	c.Result = string(jsonBytes)
}

func (c *ContextCustom) ErrorParams(msg string) {
	c.Error(406, msg)
}

func (c *ContextCustom) ErrorAuth(msg string) {
	c.Error(401, msg)
}

func (c *ContextCustom) ErrorCustom(msg string) {
	c.Error(500, msg)
}

func (c *ContextCustom) Error(code int, msg string) {
	m := gin.H{
		"code": code,
		"msg":  msg,
	}
	c.JSON(code, m)
	c.Abort()

	jsonBytes, _ := json.Marshal(m)
	c.Result = string(jsonBytes)
}
