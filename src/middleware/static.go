package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

// 静态文件
func Static(app *gin.Engine, prefix string) {
	server := app.Group(prefix)
	server.GET("/*path", func(c *gin.Context) {
		filename := prefix + c.Param("path")
		entryFilename := prefix + "/index.html"
		file, err := os.Open(filename)
		if err != nil {
			c.File(entryFilename)
			return
		}
		defer file.Close()

		info, _ := file.Stat()
		if info.IsDir() {
			c.File(entryFilename)
		} else {
			c.File(filename)
		}
	})
}
