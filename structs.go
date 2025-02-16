package main

/*
? Данные, заполняемые пользователем при авторизации и(или) регистрации.
?	Name - имя пользователя;
?	Login - псевдоним пользователя;
?	Email = электронная почта пользователя;
?	Password - пароль пользователя;
*/
type UserInput struct {
	Name     string `form:"username" binding:"required" json:"name" gorm:"column:name"`
	Login    string `form:"login" binding:"required" json:"login" gorm:"column:login"`
	Email    string `form:"email" binding:"required" json:"email" gorm:"column:email"`
	Password string `form:"password" binding:"required" json:"password" gorm:"column:password"`
}

type User struct {
	Id       uint   `gorm:"column:id"`
	Role     int8   `gorm:"column:role_code"`
	Name     string `gorm:"column:name"`
	Login    string `gorm:"column:login"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
}

type Login struct {
	LoginPrm string `form:"identifier" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type NoteType struct {
	Id     uint   `gorm:"column:id"`
	Text   string `form:"noteType" binding:"required" gorm:"column:name"`
	Author uint   `gorm:"column:creator_id"`
}
