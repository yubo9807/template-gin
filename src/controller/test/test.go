package test

import (
	"server/src/middleware"

	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	ctx := middleware.ContextGet(c)
	ctx.Success()
}
