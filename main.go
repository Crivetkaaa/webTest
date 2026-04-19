package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"webTest/database_folder"

	"github.com/gin-gonic/gin"
)

func setupStatic(r *gin.Engine) {
	r.Static("/statics", "./statics")
}

func catalogs(c *gin.Context, db *database_folder.DB) {
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

func getProducrs(c *gin.Context, db *database_folder.DB) {
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

func getProductInfo(c *gin.Context, db *database_folder.DB) {
	product := c.Param("product_slug")
	pi, err := db.GetBonusInfo(product)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	fmt.Println(pi)
	c.JSON(http.StatusOK, pi)

}

func getProduct(c *gin.Context, db *database_folder.DB) {
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

func getProductAPI(c *gin.Context, db *database_folder.DB) {
	slug := c.Param("name")

	product, err := db.GetProduct(slug)
	if err != nil {
		c.JSON(404, nil)
		return
	}

	c.JSON(200, product)
}

func setupRoutes(r *gin.Engine, db *database_folder.DB) {
	// 1. Статические роуты и группы (Сначала они!)
	api := r.Group("/api")
	{
		api.GET("/mini_categories", func(ctx *gin.Context) {
			GetMiniNavbar(ctx, db)
		})
	}

	api.GET("/get_products", func(ctx *gin.Context) {
		getProducrs(ctx, db)
	})

	api.GET("/product_info/:product_slug", func(ctx *gin.Context) {
		getProductInfo(ctx, db)
	})
	api.GET("/product/:name", func(ctx *gin.Context) {
		getProductAPI(ctx, db)
	})

	api.POST("/orders", func(ctx *gin.Context) {

	})

	r.GET("/", func(ctx *gin.Context) {
		d, err := db.GetType()
		if err != nil || d == "" {
			ctx.HTML(http.StatusNotFound, "err.html", gin.H{"Title": "Ошибка сервера"})
			return
		}
		ctx.Redirect(http.StatusMovedPermanently, "/"+d)
	})

	r.GET("/product/:name", func(ctx *gin.Context) {
		getProduct(ctx, db)
	})

	// 2. ДИНАМИЧЕСКИЙ РОУТ (В САМЫЙ КОНЕЦ)
	// Он будет срабатывать только если запрос не подошел под /api или /
	r.GET("/:type", func(ctx *gin.Context) {
		pageType := ctx.Param("type")

		// Eсли вдруг попал favicon или api (хотя порядок должен спасти)
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
		catalogs(ctx, db)
	})
}

func setupTemplates(r *gin.Engine) {
	funcMap := template.FuncMap{
		"join": strings.Join,

		"joinAny": func(v interface{}, sep string) string {
			switch t := v.(type) {
			case []string:
				return strings.Join(t, sep)
			case []int:
				arr := make([]string, len(t))
				for i, val := range t {
					arr[i] = strconv.Itoa(val)
				}
				return strings.Join(arr, sep)
			}
			return ""
		},

		"toJson": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		},
	}

	r.SetFuncMap(funcMap)
	r.LoadHTMLGlob("templates/**/*.html")
}

func main() {
	db, err := database_folder.CreateDB()
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}

	r := gin.Default()

	setupTemplates(r)
	setupStatic(r)
	setupRoutes(r, db)

	log.Println("Server started on :8080")
	r.Run(":8080")
}
