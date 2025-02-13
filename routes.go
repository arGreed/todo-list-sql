package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	pingRoute     string = "/ping"
	registerRoute string = "/register"
)

var (
	userTab string = "to_do_list.user"
)

/*
? Маршрут для проверки соединения;
*/
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"response": "pong"})
}

func showReregister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func register(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var input UserInput
		err := c.ShouldBind(&input)
		if err != nil {
			log.Println("Ошибка при получении json файла:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Получены некорректные данные"})
			return
		}
		fmt.Println(input)
		var user User

		result := db.Table(userTab).Where("login = ? or email = ?", input.Login, input.Email).First(&user)

		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			log.Println("Ошибка при работе с базой:", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при работе с базой данных"})
			return
		}

		user.Email = input.Email
		user.Login = input.Login
		user.Name = input.Name
		user.Password = input.Password

		result = db.Table(userTab).Create(&user)

		if result.Error != nil {
			log.Println("Ошибка при работе с базой:", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при работе с базой данных"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"response": "Пользователь успешно создан"})
	}
}
