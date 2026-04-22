package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Ma'lumotlar bazasi jadvali strukturasi
type Message struct {
	gorm.Model
	Content string
}

func main() {
	// Render-dagi DATABASE_URL orqali ulanish
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Ma'lumotlar bazasiga ulanib bo'lmadi!")
	}

	// Jadvalni avtomatik yaratish
	db.AutoMigrate(&Message{})

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// 1. Asosiy sahifa (Sovg'a)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Message": "Zar, men seni juda yaxshi ko'raman!",
		})
	})

	// 2. Javob yuborish funksiyasi
	r.POST("/reply", func(c *gin.Context) {
		reply := c.PostForm("reply_text")
		if reply != "" {
			db.Create(&Message{Content: reply})
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Message": "Xabaringiz yuborildi! ❤️",
		})
	})

	// 3. Admin panel (Kelgan xabarlarni ko'rish)
	r.GET("/javoblar", func(c *gin.Context) {
		var messages []Message
		db.Order("created_at desc").Find(&messages) // Yangilari tepada chiqadi
		c.HTML(http.StatusOK, "admin.html", gin.H{
			"Messages": messages,
		})
	})

	// Portni sozlash
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
