package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Blog struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CreateBlogInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

var DB *gorm.DB

func connectDatabase() {
	database, err := gorm.Open(sqlite.Open("blogo.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = database.AutoMigrate(&Blog{})
	if err != nil {
		log.Fatal(err)
	}

	DB = database
}

func getBlogs(w http.ResponseWriter, r *http.Request) {
	var blogs []Blog

	DB.Find(&blogs)

	tmpl := template.Must(template.ParseFiles("blogs.html"))
	tmpl.Execute(w, blogs)
}

func getBlog(w http.ResponseWriter, r *http.Request) {
	var blog Blog

	err := DB.Where("id = ?", r.PathValue("id")).First(&blog).Error
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(template.ParseFiles("blog.html"))
	tmpl.Execute(w, blog)
}

func createBlog(w http.ResponseWriter, r *http.Request) {
	var blog Blog

	json.NewDecoder(r.Body).Decode(&blog)

	DB.Create(&blog)
}

func main() {
	mux := http.NewServeMux()

	connectDatabase()

	mux.HandleFunc("GET /blogs", getBlogs)
	mux.HandleFunc("POST /blogs", createBlog)
	mux.HandleFunc("GET /blogs/{id}", getBlog)

	log.Print("Listening...")

	http.ListenAndServe(":8080", mux)
}
