package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	pingRoute     string = "/ping"
	registerRoute string = "/register"
	loginRoute    string = "/login"
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

func hashPassword(pass string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
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
		hash, err := hashPassword(input.Password)
		if err != nil {
			log.Println("Ошибка при хешировании пароля:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка шифрования пароля"})
			return
		}
		user.Email = input.Email
		user.Login = input.Login
		user.Name = input.Name
		user.Password = hash

		result = db.Table(userTab).Create(&user)

		if result.Error != nil {
			log.Println("Ошибка при работе с базой:", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при работе с базой данных"})
			return
		}

		c.Redirect(http.StatusSeeOther, "/")
	}
}

func showIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func showLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func login(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var login Login

		err := c.ShouldBind(&login)
		if err != nil {
			log.Println("Получены некорректные данные:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Получены некорректные данные"})
			return
		}
		var user User
		result := db.Table(userTab).Where("login = ? or email = ?", login.LoginPrm, login.LoginPrm).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				log.Println("Ошибка работы с бд")
				c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
				return
			} else {
				log.Println("Ошибка работы с бд:", result.Error)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка работы с бд"})
				return
			}
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
		if err != nil {
			log.Println("Пароли не совпадают:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пароли не совпадают"})
			return
		}
		c.Redirect(http.StatusSeeOther, "/")
	}
}
