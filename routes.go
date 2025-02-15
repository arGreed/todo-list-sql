package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	pingRoute        string = "/ping"
	registerRoute    string = "/register"
	loginRoute       string = "/login"
	addNoteTypeRoute string = "/note/add"
	logoutRoute      string = "/logout"
)

var (
	userTab string = "to_do_list.user"
)

var jwtSecret = []byte("ergoipahmjn-weomfwep4oghjmethomer[gp]")

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
	userId, exists := c.Get("userId")
	isAuthenticated := exists && userId != nil
	log.Printf("Рендеринг главной страницы. Auth: %v, UserID: %v", isAuthenticated, userId)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"isAuthenticated": isAuthenticated,
	})
}

func showLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func generateToken(user uint, role int8) string {
	claims := jwt.MapClaims{}

	claims["authorized"] = true
	claims["userId"] = user
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return ""
	}
	return tokenString
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			log.Println("Кука auth_token отсутствует")
			c.Next()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неверный метод подписи")
			}
			return jwtSecret, nil
		})

		if err != nil {
			log.Println("Ошибка парсинга токена:", err)
			c.Next()
			return
		}

		if !token.Valid {
			log.Println("Токен невалиден")
			c.Next()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			log.Printf("Успешная аутентификация. UserID: %v", claims["userId"])
			c.Set("userId", claims["userId"])
			c.Set("role", claims["role"])
		}

		c.Next()
	}
}

func StrictAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("userId"); !exists {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
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

		token := generateToken(user.Id, user.Role)
		if token == "" {
			log.Println("Ошибка генерации токена")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
			return
		}

		c.SetCookie(
			"auth_token",
			token,
			3600*24,
			"/",
			"",
			false,
			true,
		)

		c.Redirect(http.StatusSeeOther, "/")
	}
}

func logout(c *gin.Context) {
	c.SetCookie(
		"auth_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.Redirect(http.StatusSeeOther, "/")
}

func showAddNoteType(c *gin.Context) {
	c.HTML(http.StatusOK, "addNoteType.html", nil)
}

//func addNoteType(db *gorm.DB) func(c *gin.Context) {
//	return func(c *gin.Context)
//}
