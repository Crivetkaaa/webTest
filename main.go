package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/big" // Добавили обратно для стандартного веб-сервера
	"strconv"
	"strings"

	"webTest/database_folder"
	"webTest/middleware"
	"webTest/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupStatic(r *gin.Engine) {
	// ИЗМЕНЕНО ДЛЯ VDS: Возвращаем локальный путь к статике
	r.Static("/statics", "./statics")
}

func setupRoutes(r *gin.Engine, db *database_folder.DB) {
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
	// ИЗМЕНЕНО ДЛЯ VDS: Возвращаем локальный путь к шаблонам
	r.LoadHTMLGlob("templates/**/*.html")
}

func GenerateRandomPassword(length int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result[i] = letters[num.Int64()]
	}
	return string(result), nil
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	err := godotenv.Load()
	if err != nil {
		log.Println("Предупреждение: .env файл не найден, берутся системные переменные")
	}

	db, err := database_folder.CreateDB()
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}

	var count int
	err = db.Db.QueryRow("SELECT COUNT(*) FROM administrator").Scan(&count)
	if err != nil {
		panic("Failed to check database state: " + err.Error())
	}

	if count == 0 {
		login, err := GenerateRandomPassword(20)
		if err != nil {
			log.Fatalf("Login error: %s", err)
		}
		password, err := GenerateRandomPassword(20)
		if err != nil {
			log.Fatalf("Password error: %s", err)
		}
		db.Registration(login, password)
		// ИЗМЕНЕНО ДЛЯ VDS: Теперь этот вывод вы увидите прямо в консоли при запуске приложения!
		fmt.Printf("Initial admin created!\nLogin: %s\nPassword: %s\n", login, password)
	}
	middleware.InitMiddleware()

	r := gin.Default()

	setupTemplates(r)
	setupStatic(r)
	setupRoutes(r, db)

	log.Println("Server started on :8080")
	err = r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
