package routers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"webTest/database_folder"
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
		ctx.JSON(400, gin.H{"error": "Invalid data"})
		return
	}

	fmt.Println(data)

	if db.Login(data.UserName, data.UserPassword) {
		sessionID, err := generateSessionID()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Server Error"})
			return
		}

		ctx.SetCookie(
			"sessionID",
			sessionID,
			3600,
			"/admin",
			"",
			false,
			true,
		)

		ctx.JSON(http.StatusOK, gin.H{"success": true, "redirect": "/admin/dashboard", "sessionID": sessionID})
		middleware.SessionID = sessionID
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Wrong credentials"})
	}
}

func dashboardInfo(ctx *gin.Context, db *database_folder.DB) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	status := ctx.DefaultQuery("status", "all")

	orders, err := db.GetOrdersAdmin(limit, offset, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		ctx.HTML(200, "admin_login.html", nil)
	})

	admin.POST("/auth", func(ctx *gin.Context) {
		adminLogin(ctx, db)
	})

	adminPanel := admin.Group("/")
	adminPanel.Use(middleware.AuthMiddleware())

	adminPanel.GET("/dashboard", func(ctx *gin.Context) {
		ctx.HTML(200, "admin.html", nil)
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
}
