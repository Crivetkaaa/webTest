package routers

import (
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
		ctx.Redirect(http.StatusMovedPermanently, "/"+d)
	})

	r.GET("/product/:name", func(ctx *gin.Context) {
		handlers.GetProduct(ctx, db)
	})

	r.GET("/:type/:minicat", func(ctx *gin.Context) {
		t := ctx.Param("type")
		cat, err := db.GetTypes()
		if err != nil {
			ctx.HTML(404, "err.html", gin.H{"Title": "1"})
			return
		}
		if !slices.Contains(cat, t) {
			ctx.HTML(404, "err.html", gin.H{"Title": "2"})
			return
		}

		miniCat, err := db.GetMiniTypes()

		if !slices.Contains(miniCat, ctx.Param("minicat")) {
			ctx.HTML(404, "err.html", gin.H{"Title": "3"})
			return
		}

		handlers.Catalogs(ctx, db)
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

	r.GET("/download/:type", func(ctx *gin.Context) {
		fileType := ctx.Param("type")

		var filePath string
		if fileType == "privacy" {
			filePath = "statics/docx/privacy.docx"
			fileName := "privacy_policy.docx"
			handlers.Download(ctx, filePath, fileName)
		} else if fileType == "offer" {
			filePath = "statics/docx/offer.docx"
			fileName := "offer_policy.docx"
			handlers.Download(ctx, filePath, fileName)
		} else {
			ctx.HTML(404, "err.html", nil)
		}
	})
}
