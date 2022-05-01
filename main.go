package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	Username   string `json:"username" gorm:"unique"`
	Password   string `json:"password"`
}

var db *gorm.DB = getDB("user")

func main() {
	router := getRouter()
	router.Run()
}

func getRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/users/", getUsersHandler)
	router.POST("/users/", postUsersHandler)
	router.GET("/users/:username", getUserHandler)
	router.PUT("/users/:username", putUserHandler)
	router.DELETE("/users/:username", deleteUserHandler)
	return router
}

func putUserHandler(c *gin.Context) {
	username := c.Param("username")
	var user User
	if result := db.Where(&User{Username: username}).First(&user); result.Error != nil {
		c.JSON(400, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	var reqBody User
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	if reqBody.Username != "" {
		user.Username = reqBody.Username
	}
	if reqBody.Password != "" {
		user.Password = reqBody.Password
	}
	if result := db.Save(&user); result.Error != nil {
		c.JSON(400, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	c.JSON(200, user)
}

func deleteUserHandler(c *gin.Context) {
	username := c.Param("username")
	var user User
	if result := db.Where(&User{Username: username}).First(&user).Delete(&user); result.Error != nil {
		c.JSON(400, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": username + " deleted",
	})
}

func getUserHandler(c *gin.Context) {
	username := c.Param("username")
	var user User
	if result := db.Where(&User{Username: username}).First(&user); result.Error != nil {
		c.JSON(400, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	c.JSON(200, user)
}

func postUsersHandler(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	if result := db.Create(&user); result.Error != nil {
		c.JSON(400, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	c.JSON(200, user)
}

func getUsersHandler(c *gin.Context) {
	var users []User
	if result := db.Find(&users); result.Error != nil {
		c.JSON(400, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	c.JSON(200, users)
}

func getDB(dbName string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		panic("migrate error")
	}
	return db
}
