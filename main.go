package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Candidate struct {
	FullName   string
	Phone      string
	Age        string
	Address    string
	Experience string
	JobTitle   string
}

func main() {
	// 1. Ma'lumotlar bazasiga ulanish
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Jadvalni yaratish
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS candidates (
		id SERIAL PRIMARY KEY,
		full_name TEXT,
		phone TEXT,
		age TEXT,
		address TEXT,
		experience TEXT,
		job_title TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// 2. Statik fayllar (Rasm ko'rinishi uchun)
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Bosh sahifa
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Anketa yuborish
	r.POST("/apply", func(c *gin.Context) {
		var cand Candidate
		cand.FullName = c.PostForm("full_name")
		cand.Phone = c.PostForm("phone")
		cand.Age = c.PostForm("age")
		cand.Address = c.PostForm("address")
		cand.Experience = c.PostForm("experience")
		cand.JobTitle = c.PostForm("job_title")

		// Bazaga saqlash
		query := `INSERT INTO candidates (full_name, phone, age, address, experience, job_title) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err := db.Exec(query, cand.FullName, cand.Phone, cand.Age, cand.Address, cand.Experience, cand.JobTitle)
		if err != nil {
			c.String(http.StatusInternalServerError, "Xatolik: Bazaga saqlanmadi")
			return
		}

		// Telegramga yuborish (TOKEN va ID'ni o'zingiznikiga almashtiring)
		botToken := "8256517822:AAH1tYGXqxWa9w1OwtWaL4byr4oy5jzuxik"
		chatID := "5531320866"

		msgText := fmt.Sprintf(" YANGI ANKETA QABUL QILINDI:\n\n"+
			" F.I.SH: %s\n"+
			"Telefon nomer: %s\n"+
			" Yosh: %s\n"+
			"Manzil: %s\n"+
			" Tajriba: %s\n"+
			" Lavozim: %s",
			cand.FullName, cand.Phone, cand.Age, cand.Address, cand.Experience, cand.JobTitle)

		telegramURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s",
			botToken, chatID, url.QueryEscape(msgText))

		http.Get(telegramURL)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Success": "Ma'lumotlaringiz muvaffaqiyatli yuborildi!",
		})
	})

	// Render uchun port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
