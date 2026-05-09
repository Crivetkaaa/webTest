package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

func DeleteProduct(c *gin.Context, db *database_folder.DB) {
	fmt.Println("Del")
	product_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "incorrect product id"})
		return
	}

	if err := db.DeleteProduct(product_id); err != nil {
		c.JSON(400, gin.H{"error": "database error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "good"})
}

func GetProductInfo(c *gin.Context, db *database_folder.DB) {
	product := c.Param("product_slug")
	pi, err := db.GetBonusInfo(product)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "scsdf"})
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

func UpdateProduct(c *gin.Context, db *database_folder.DB) {
	// ======================
	// 1. ID
	// ======================
	idStr := c.PostForm("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt <= 0 {
		c.JSON(400, gin.H{"error": "invalid product id"})
		return
	}

	var data struct_folder.UpdateProductData
	data.ID = idInt
	data.Name = strings.TrimSpace(c.PostForm("name"))
	data.Description = strings.TrimSpace(c.PostForm("description"))

	// ======================
	// 2. ВАЛИДАЦИЯ
	// ======================
	if data.Name == "" {
		c.JSON(400, gin.H{"error": "name is required"})
		return
	}

	// ======================
	// 3. JSON
	// ======================
	if err := json.Unmarshal([]byte(c.PostForm("existingPhotos")), &data.ExistingPhotos); err != nil {
		c.JSON(400, gin.H{"error": "invalid existingPhotos"})
		return
	}

	if err := json.Unmarshal([]byte(c.PostForm("variants")), &data.Variants); err != nil {
		c.JSON(400, gin.H{"error": "invalid variants"})
		return
	}

	if err := json.Unmarshal([]byte(c.PostForm("characteristics")), &data.Characteristics); err != nil {
		c.JSON(400, gin.H{"error": "invalid characteristics"})
		return
	}

	if err := json.Unmarshal([]byte(c.PostForm("subcategories")), &data.Subcategories); err != nil {
		c.JSON(400, gin.H{"error": "invalid subcategories"})
		return
	}

	// доп проверка variants
	if len(data.Variants.Value) != len(data.Variants.Price) {
		c.JSON(400, gin.H{"error": "variants mismatch"})
		return
	}

	// ======================
	// 4. ФАЙЛЫ (временно)
	// ======================
	var tempFiles []string

	form, err := c.MultipartForm()
	if err == nil && form.File != nil {
		files := form.File["newPhotos"]

		for _, file := range files {
			fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
			dst := "statics/img/product_img/" + fileName

			if err := c.SaveUploadedFile(file, dst); err != nil {
				c.JSON(500, gin.H{"error": "failed to save file"})
				return
			}

			tempFiles = append(tempFiles, dst)
			data.NewPhotoPaths = append(data.NewPhotoPaths, dst)
		}
	}

	// ======================
	// 5. DB
	// ======================
	err = db.UpdateProduct(data)
	if err != nil {
		// 🔴 если БД упала — удаляем загруженные файлы
		for _, f := range tempFiles {
			os.Remove(f)
		}

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// ======================
	// OK
	// ======================
	c.JSON(200, gin.H{"status": "success"})
}

func AddProduct(c *gin.Context, db *database_folder.DB) {
	var data struct_folder.UpdateProductData

	data.Name = strings.TrimSpace(c.PostForm("name"))
	data.Description = strings.TrimSpace(c.PostForm("description"))

	if data.Name == "" {
		c.JSON(400, gin.H{"error": "name is required"})
		return
	}

	// JSON
	if err := json.Unmarshal([]byte(c.PostForm("variants")), &data.Variants); err != nil {
		c.JSON(400, gin.H{"error": "invalid variants"})
		return
	}

	if err := json.Unmarshal([]byte(c.PostForm("characteristics")), &data.Characteristics); err != nil {
		c.JSON(400, gin.H{"error": "invalid characteristics"})
		return
	}

	if err := json.Unmarshal([]byte(c.PostForm("subcategories")), &data.Subcategories); err != nil {
		c.JSON(400, gin.H{"error": "invalid subcategories"})
		return
	}

	// файлы
	form, err := c.MultipartForm()
	if err == nil && form.File != nil {
		files := form.File["newPhotos"]

		for _, file := range files {
			fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
			dst := "statics/img/product_img/" + fileName

			if err := c.SaveUploadedFile(file, dst); err != nil {
				c.JSON(500, gin.H{"error": "file upload error"})
				return
			}

			data.NewPhotoPaths = append(data.NewPhotoPaths, dst)
		}
	}

	// DB
	id, err := db.CreateProduct(data)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"id":     id,
	})
}

type CreateCategoryRequest struct {
	Name   string `json:"name"`
	Parent string `json:"parent_slug"` // slug родительской категории
}

func CreateCategory(c *gin.Context, db *database_folder.DB) {
	var req CreateCategoryRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid json",
		})
		return
	}
	req.Name = strings.TrimSpace(req.Name)

	if req.Name == "" {
		c.JSON(400, gin.H{
			"error": "empty name",
		})
		return
	}

	if req.Parent != "" {
		if err := db.AddSubcategory(req.Name, req.Parent); err != nil {
			c.JSON(400, gin.H{
				"error": "database error",
			})
			return
		}
	} else {
		if err := db.AddCategory(req.Name); err != nil {
			c.JSON(400, gin.H{
				"error": "database error",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "good"})
}

type DeleteCategoryStruct struct {
	Slug  string `json:"slug"`
	Types string `json:"type"`
}

func DeleteCategory(c *gin.Context, db *database_folder.DB) {
	var res DeleteCategoryStruct
	err := c.BindJSON(&res)

	if err != nil {
		c.JSON(400, gin.H{"error": "invalid JSON"})
		return
	}
	var table string

	switch res.Types {
	case "category":
		table = "categories"
	case "subcategory":
		table = "subcategories"
	default:
		c.JSON(400, gin.H{"error": "error category type"})
		return

	}

	if err := db.DeleteCategory(res.Slug, table); err != nil {
		c.JSON(400, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "good"})
}

type UpdateCat struct {
	CatNewName string `json:"name"`
	Types      string `json:"type"`
	Slug       string `json:"slug"`
}

func UpdateCategory(c *gin.Context, db *database_folder.DB) {
	var res UpdateCat

	err := c.BindJSON(&res)
	if err != nil {
		c.JSON(400, gin.H{"error": "incorect JSON"})
		return
	}
	var table string
	switch res.Types {
	case "category":
		table = "categories"
	case "subcategory":
		table = "subcategories"
	default:
		c.JSON(400, gin.H{"error": "error category type"})
		return
	}

	if err := db.UpdateCategory(res.CatNewName, res.Slug, table); err != nil {
		c.JSON(400, gin.H{"error": "database error"})
		return
	}

	c.JSON(200, gin.H{"message": "good"})

}
