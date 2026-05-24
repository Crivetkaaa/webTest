package routers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strconv"

	"webTest/database_folder"
	"webTest/handlers"
	"webTest/middleware"

	"github.com/gin-gonic/gin"
)

func generateSessionID() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func adminLogin(ctx *gin.Context, db *database_folder.DB) {
	type Admin struct {
		UserName     string `json:"username"`
		UserPassword string `json:"password"`
	}

	var data Admin

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data",
		})
		return
	}

	if !db.Login(data.UserName, data.UserPassword) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Wrong credentials",
		})
		return
	}

	sessionID, err := generateSessionID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Server error",
		})
		return
	}

	middleware.AddSession(sessionID, data.UserName)

	ctx.SetCookie(
		"sessionID",
		sessionID,
		3600,
		"/admin",
		"",
		true, // secure=true в production
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{
		"success":  true,
		"redirect": "/admin/dashboard",
	})
}

func adminLogout(ctx *gin.Context) {
	sessionID, err := ctx.Cookie("sessionID")

	if err == nil && sessionID != "" {
		middleware.DeleteSession(sessionID)
	}

	ctx.SetCookie(
		"sessionID",
		"",
		-1,
		"/admin",
		"",
		true,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func dashboardInfo(ctx *gin.Context, db *database_folder.DB) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	status := ctx.DefaultQuery("status", "all")

	orders, err := db.GetOrdersAdmin(limit, offset, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func updateStatus(ctx *gin.Context, db *database_folder.DB) {
	type NewStatus struct {
		Id        int    `json:"id"`
		NewStatus string `json:"status"`
	}

	var status NewStatus

	err := ctx.ShouldBindJSON(&status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	err = db.UpdateStatus(status.Id, status.NewStatus)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func AdminRouters(admin *gin.RouterGroup, db *database_folder.DB) {
	admin.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "admin_login.html", nil)
	})

	admin.POST("/auth", func(ctx *gin.Context) {
		adminLogin(ctx, db)
	})

	adminPanel := admin.Group("/")
	adminPanel.Use(middleware.AuthMiddleware())

	adminPanel.POST("/logout", func(ctx *gin.Context) {
		adminLogout(ctx)
	})

	adminPanel.GET("/dashboard", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "admin.html", nil)
	})

	adminPanel.GET("/products", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminProduct.html", nil)
	})

	adminPanel.GET("/dashboard/info", func(ctx *gin.Context) {
		dashboardInfo(ctx, db)
	})

	adminPanel.POST("/orders/status", func(ctx *gin.Context) {
		updateStatus(ctx, db)
	})

	adminPanel.POST("/update_product", func(ctx *gin.Context) {
		handlers.UpdateProduct(ctx, db)
	})

	adminPanel.POST("/create_product", func(ctx *gin.Context) {
		handlers.AddProduct(ctx, db)
	})

	adminPanel.DELETE("/delete_product/:id", func(ctx *gin.Context) {
		handlers.DeleteProduct(ctx, db)
	})

	adminPanel.POST("/create_category", func(ctx *gin.Context) {
		handlers.CreateCategory(ctx, db)
	})

	adminPanel.DELETE("/delete_category", func(ctx *gin.Context) {
		handlers.DeleteCategory(ctx, db)
	})

	adminPanel.POST("/update_category", func(ctx *gin.Context) {
		handlers.UpdateCategory(ctx, db)
	})

	adminPanel.GET("/settings", func(ctx *gin.Context) {
		handlers.AdminSettings(ctx)
	})

	adminPanel.POST("/settings/update-docs", func(ctx *gin.Context) {
		handlers.AdminUpdateDocx(ctx)
	})

	adminPanel.POST("/settings/change-password", func(ctx *gin.Context) {
		handlers.ChangePassword(ctx, db)
	})
}
