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

// Rezyume ma'lumotlari uchun model
type Candidate struct {
	gorm.Model
	FullName   string
	Phone      string
	Age        string
	Address    string
	Experience string
	JobTitle   string
}

func main() {
	// Ma'lumotlar bazasiga ulanish
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Ma'lumotlar bazasiga ulanishda xato!")
	}
	db.AutoMigrate(&Candidate{})

	r := gin.Default()

	// Statik fayllar (rasm, css) va shablonlar
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Asosiy sahifa
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Anketani qabul qilish
	r.POST("/apply", func(c *gin.Context) {
		cand := Candidate{
			FullName:   c.PostForm("full_name"),
			Phone:      c.PostForm("phone"),
			Age:        c.PostForm("age"),
			Address:    c.PostForm("address"),
			Experience: c.PostForm("experience"),
			JobTitle:   c.PostForm("job_title"),
		}

		db.Create(&cand)

		// TELEGRAM BILDIRISHNOMASI
		// Bu yerga o'z ma'lumotlaringizni qo'ying
		botToken := "SIZNING_BOT_TOKENINGIZ"
		chatID := "SIZNING_ID_RAQAMINGIZ"

		msgText := fmt.Sprintf(" YANGI REZYUME QABUL QILINDI:\n\n"+
			" F.I.SH: %s\n"+
			" Telefon Manzili: %s\n"+
			" Ish tajribasi: %s\n"+
			"Qaysi Lavozimda ishlash xoxishi: %s",
			cand.FullName, cand.Phone, cand.Age, cand.Address, cand.Experience, cand.JobTitle)

		telegramURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s",
			botToken, chatID, url.QueryEscape(msgText))

		// Telegramga yuborish
		http.Get(telegramURL)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Success": "Ma'lumotlaringiz muvaffaqiyatli yuborildi! Tez orada siz bilan bog'lanamiz.",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
