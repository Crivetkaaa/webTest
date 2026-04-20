package main

import (
	"encoding/json"
	"html/template"
	"log"
	"strconv"
	"strings"

	"webTest/database_folder"
	"webTest/routers"

	"github.com/gin-gonic/gin"
)

func setupStatic(r *gin.Engine) {
	r.Static("/statics", "./statics")
}

func setupRoutes(r *gin.Engine, db *database_folder.DB) {
	// 1. Статические роуты и группы (Сначала они!)
	api := r.Group("/api")
	admin := r.Group("/admin")
	users := r.Group("/")
	routers.APIRouters(api, db)
	routers.AdminRouters(admin, db)
	routers.UsersRouters(users, db)

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
