package routers

import (
	"webTest/database_folder"
	"webTest/handlers"

	"github.com/gin-gonic/gin"
)

func APIRouters(api *gin.RouterGroup, db *database_folder.DB) {
	api.GET("/mini_categories", func(ctx *gin.Context) {
		handlers.GetMiniNavbar(ctx, db)
	})

	api.GET("/get_products", func(ctx *gin.Context) {
		handlers.GetProducrs(ctx, db)
	})

	api.GET("/product_info/:product_slug", func(ctx *gin.Context) {
		handlers.GetProductInfo(ctx, db)
	})
	api.GET("/product/:name", func(ctx *gin.Context) {
		handlers.GetProductAPI(ctx, db)
	})

	api.POST("/orders", func(ctx *gin.Context) {
		handlers.PostOrders(ctx, db)
	})

	api.GET("/categories", func(ctx *gin.Context) {
		handlers.GetCategories(ctx, db)
	})
}
