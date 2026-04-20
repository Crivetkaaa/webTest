package routers

import (
	"webTest/database_folder"

	"github.com/gin-gonic/gin"
)

func AdminRouters(admin *gin.RouterGroup, db *database_folder.DB) {
	admin.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "admin_login.html", nil)
	})

	admin.POST("/auth", func(ctx *gin.Context) {
		type Admin struct {
			UserName     string `json:"username"`
			UserPassword string `json:"password"`
		}
		var data Admin

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid data"})
			return
		}

		if data.UserName == "admin" && data.UserPassword == "1234" {
			// Возвращаем успех и URL, куда нужно перейти
			ctx.JSON(200, gin.H{"success": true, "redirect": "/admin/dashboard"})
		} else {
			ctx.JSON(401, gin.H{"success": false, "message": "Wrong credentials"})
		}
	})
}
