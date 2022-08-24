package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

type Connect struct {
	Name   string
	Number string
}

var Con Connect

func ReturnStatus(r *gin.Engine, c *gin.Context, err error) {
	if err != nil {
		c.Request.URL.Path = "/status/failed"
	}
	c.Request.URL.Path = "/status/success"
	r.HandleContext(c)
}

func RouteInit(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	cre := r.Group("create")
	{
		cre.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "create.html", nil)
		})
		cre.POST("/", func(c *gin.Context) {
			Con.Name = c.PostForm("name")
			Con.Number = c.PostForm("number")
			err := db.Create(&Con).Error
			ReturnStatus(r, c, err)
		})
	}
	del := r.Group("delete")
	{
		del.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "delete.html", nil)
		})
		del.POST("/", func(c *gin.Context) {
			Con.Name = c.PostForm("name")
			err := db.Where("name = ?", Con.Name).Delete(&Connect{}).Error
			ReturnStatus(r, c, err)
		})
	}
	upd := r.Group("update")
	{
		upd.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "update.html", nil)
		})
		upd.POST("/", func(c *gin.Context) {
			Con.Name = c.PostForm("name")
			Con.Number = c.PostForm("number")
			err := db.Create(&Con).Error
			ReturnStatus(r, c, err)
		})
	}
	read := r.Group("read")
	{
		read.GET("/", func(c *gin.Context) {
			contacts := Read(db)
			c.HTML(http.StatusOK, "read.tmpl", gin.H{
				"res": contacts,
			})
		})
	}
	stt := r.Group("status")
	{
		stt.POST("/success", func(c *gin.Context) {
			c.HTML(http.StatusOK, "success.html", nil)
		})
		stt.POST("/failed", func(c *gin.Context) {
			c.HTML(http.StatusOK, "failed.html", nil)
		})
	}
	return r
}

func ConnectDB() *gorm.DB {
	dsn := "root:Mikasa521@tcp(127.0.0.1:3306)/address_book?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		fmt.Println("Connect database err:", err)
	}
	err = db.AutoMigrate(&Connect{})
	if err != nil {
		fmt.Println("Migrate err", err)
	}
	return db
}

func Read(db *gorm.DB) []Connect {
	var connects []Connect
	db.Find(&connects)
	fmt.Println(connects)
	return connects
}

func main() {
	db := ConnectDB()
	r := RouteInit(db)

	err := r.Run()
	if err != nil {
		fmt.Println("r.Run err:", err)
	}
}
