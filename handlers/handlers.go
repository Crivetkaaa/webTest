package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"webTest/database_folder"
	"webTest/struct_folder"

	"github.com/gin-gonic/gin"
)

func Catalogs(c *gin.Context, db *database_folder.DB) {
	bigNavbar, err := db.GetBigNavbar()
	if err != nil {
		c.String(300, err.Error())
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title":     "Каталог",
		"Product":   nil,
		"BigNavbar": bigNavbar,
	})
}

func GetMiniNavbar(c *gin.Context, db *database_folder.DB) {
	cat := c.Query("category")
	navbar, err := db.GetMiniNavbar(cat)
	fmt.Println(navbar)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	c.JSON(http.StatusOK, navbar)
}

func GetProducrs(c *gin.Context, db *database_folder.DB) {
	productType := c.Query("type")
	productCategory := c.Query("category")

	limitStr := c.DefaultQuery("limit", "20")  // Если пусто, будет 20
	offsetStr := c.DefaultQuery("offset", "0") // Если пусто, будет 0

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	products, err := db.GetProducts(productType, productCategory, offset, limit)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	fmt.Println(len(products))
	c.JSON(http.StatusOK, products)
}

func GetProductInfo(c *gin.Context, db *database_folder.DB) {
	product := c.Param("product_slug")
	pi, err := db.GetBonusInfo(product)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	fmt.Println(pi)
	c.JSON(http.StatusOK, pi)

}

func GetProduct(c *gin.Context, db *database_folder.DB) {
	productSlug := c.Param("name")

	product_info, err := db.GetProduct(productSlug)
	if err != nil {
		c.HTML(http.StatusNotFound, "err.html", gin.H{"Title": "Что-то пошло не так"})
		return
	}

	c.HTML(http.StatusOK, "product_page.html", gin.H{
		"Product": product_info,
		"Title":   "Каталог",
	})
}

func GetProductAPI(c *gin.Context, db *database_folder.DB) {
	slug := c.Param("name")

	product, err := db.GetProduct(slug)
	if err != nil {
		c.JSON(404, nil)
		return
	}
	c.JSON(200, product)
}

func PostOrders(c *gin.Context, db *database_folder.DB) {
	var order struct_folder.OrderData

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(order)

	err := db.InsertOrder(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при сохранении заказа: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func GetCategories(c *gin.Context, db *database_folder.DB) {
	nav, err := db.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	c.JSON(http.StatusOK, nav)

}
