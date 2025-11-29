package middleware

import (
	"server/src/service"

	"github.com/gin-gonic/gin"
)

const KEY = "user_info"

func Authorization(c *gin.Context) {
	ctx := ContextGet(c)

	auth := ctx.GetHeader("Authorization")
	if auth == "" {
		ctx.ErrorAuth("Unauthorized")
		ctx.Abort()
		return
	}

	info, err := service.Jwt.Verify(auth)
	if err != nil {
		if err.Error() == "Token is expired" {
			// 可尝试刷新 token
			ctx.ErrorAuth(err.Error())
		} else {
			ctx.ErrorAuth(err.Error())
		}
		ctx.Abort()
		return
	}
	ctx.Set(KEY, info)
}

// 获取 token 储存信息
func GetTokenInfo(ctx *gin.Context) map[string]interface{} {
	info, _ := ctx.Get(KEY)
	return info.(map[string]interface{})
}

// 角色校验
func RoleVerify(roleId string) gin.HandlerFunc {
	return func(c *gin.Context) {
		info := GetTokenInfo(c)
		ctx := ContextGet(c)
		if info["roleId"] != roleId {
			ctx.ErrorCustom("The current user has no permission")
			return
		}
	}
}
