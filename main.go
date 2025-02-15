package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dsn string = "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable TimeZone=UTC"

func dbInit() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	router := gin.Default()
	db, err := dbInit()
	router.LoadHTMLGlob("templates/*")
	if err != nil {
		log.Println("Ошибка инициализации базы данных:", err)
		return
	}
	router.Use(AuthMiddleware())

	router.GET(pingRoute, ping)
	router.GET("/", showIndex)
	router.GET(registerRoute, showReregister)
	router.GET(loginRoute, showLogin)
	router.POST(loginRoute, login(db))
	router.POST(registerRoute, register(db))
	router.GET(logoutRoute, logout)

	protected := router.Group("/")
	protected.Use(StrictAuthMiddleware())
	{
		protected.GET(addNoteTypeRoute, showAddNoteType)
	}

	router.Run(":8081")
}
