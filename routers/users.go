package routers

import (
	"fmt"
	"net/http"
	"slices"
	"webTest/database_folder"
	"webTest/handlers"

	"github.com/gin-gonic/gin"
)

func UsersRouters(r *gin.RouterGroup, db *database_folder.DB) {
	r.GET("/", func(ctx *gin.Context) {
		d, err := db.GetType()
		if err != nil || d == "" {
			ctx.HTML(http.StatusNotFound, "err.html", gin.H{"Title": "Ошибка сервера"})
			return
		}
		fmt.Println(d)
		ctx.Redirect(http.StatusMovedPermanently, "/"+d)
	})

	r.GET("/product/:name", func(ctx *gin.Context) {
		handlers.GetProduct(ctx, db)
	})

	// 2. ДИНАМИЧЕСКИЙ РОУТ (В САМЫЙ КОНЕЦ)
	// Он будет срабатывать только если запрос не подошел под /api или /
	r.GET("/:type", func(ctx *gin.Context) {
		pageType := ctx.Param("type")

		if pageType == "api" || pageType == "statics" || pageType == "favicon.ico" {
			ctx.Abort() // Не обрабатываем здесь
			return
		}

		categories, err := db.GetTypes()
		if err != nil {
			ctx.HTML(http.StatusNotFound, "err.html", gin.H{"Title": "Ошибка сервера"})
			return
		}

		if !slices.Contains(categories, pageType) {
			ctx.HTML(http.StatusNotFound, "err.html", gin.H{"Title": "Ошибка адреса"})
			return
		}
		handlers.Catalogs(ctx, db)
	})
}
