package user

import (
	"server/src/middleware"
	"server/src/service"

	"github.com/gin-gonic/gin"
)

// 登录
func SignIn(c *gin.Context) {
	ctx := middleware.ContextGet(c)

	type Params struct {
		Username string `binding:"required"`
		Password string `binding:"required"`
	}
	var params Params
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.ErrorParams(err.Error())
		return
	}

	// 验证用户名密码

	// 将一些重要信息存在 token 中
	info := map[string]interface{}{
		"roleId":   "0",
		"userId":   "0",
		"username": params.Username,
	}
	token := service.Jwt.Publish(info)
	ctx.SuccessData(token)
}
