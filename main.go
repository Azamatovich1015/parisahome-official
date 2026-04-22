package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Content string
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Baza bilan bog'lanishda xato!")
	}
	db.AutoMigrate(&Message{})

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// Asosiy sahifa
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Message": "Zar, men seni juda yaxshi ko'raman!",
		})
	})

	// Javob yuborish va Telegram bildirishnomasi
	r.POST("/reply", func(c *gin.Context) {
		reply := c.PostForm("reply_text")
		if reply != "" {
			db.Create(&Message{Content: reply})

			// TELEGRAM SOZLAMALARI (O'z ma'lumotlaringizni qo'ying)
			botToken := "SIZNING_BOT_TOKENINGIZ"
			chatID := "SIZNING_ID_RAQAMINGIZ"
			text := fmt.Sprintf("Yangi javob keldi: %s", reply)

			telegramURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s",
				botToken, chatID, url.QueryEscape(text))
			http.Get(telegramURL)
		}
		c.HTML(http.StatusOK, "index.html", gin.H{"Message": "Xabaringiz yuborildi! ❤️"})
	})

	// Admin panel (Parol: admin123)
	r.GET("/javoblar", func(c *gin.Context) {
		pass := c.Query("pw")
		if pass != "20021015" {
			c.String(http.StatusForbidden, "Kirish taqiqlangan! Parol noto'g'ri.")
			return
		}

		var messages []Message
		db.Order("created_at desc").Find(&messages)
		c.HTML(http.StatusOK, "admin.html", gin.H{"Messages": messages})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
