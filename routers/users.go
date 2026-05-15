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

	r.GET("/download/privacy", func(ctx *gin.Context) {
		filePath := "statics/docx/privacy.docx"

		// Задаем заголовок, который заставит браузер именно СКАЧАТЬ файл, а не пытаться открыть как текст
		ctx.Header("Content-Description", "File Transfer")
		ctx.Header("Content-Disposition", "attachment; filename=privacy_policy.docx")
		ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
		ctx.Header("Content-Transfer-Encoding", "binary")
		ctx.Header("Expires", "0")
		ctx.Header("Cache-Control", "must-revalidate")
		ctx.Header("Pragma", "public")

		// Отдаем файл в поток
		ctx.File(filePath)
	})
}
