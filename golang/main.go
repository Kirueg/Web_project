package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Year        int    `json:"year"`
	Description string `json:"description"`
	Actors      string `json:"actors"`
	Directors   string `json:"directors"`
	Subtitles   string `json:"subtitles"`
	ImagePath   string `json:"image_path"`
}

var db *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "89858243234"
	dbname   = "project"
)

func main() {
	// Создаем папку для загрузки файлов, если она не существует
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		fmt.Println("Ошибка создания папки uploads:", err)
		return
	}

	// Подключение к базе данных
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer db.Close()

	// Проверка подключения к базе данных
	if err := db.Ping(); err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}

	// Создание таблицы, если она не существует
	createTable()

	// Настройка маршрутов
	r := gin.Default()

	// Подключаем CORS-мидлвар
	r.Use(corsMiddleware())

	// Маршрут для статических файлов (изображений)
	r.Static("/uploads", "./uploads")

	r.GET("/api/movies/:id", getMovieByID)
	r.POST("/api/add-movie", handleAddMovie)
	r.GET("/api/movies", handleGetMovies)

	// Запуск сервера
	fmt.Println("Сервер запущен на http://localhost:8080")
	r.Run(":8080")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func createTable() {
	query := `
		CREATE TABLE IF NOT EXISTS movies (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			year INT NOT NULL,
			description TEXT,
			actors TEXT,
			directors TEXT,
			subtitles TEXT,
			image_path TEXT
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Ошибка создания таблицы:", err)
	}
}

func getMovieByID(c *gin.Context) {
	// Получаем ID из URL
	id := c.Param("id")

	// Запрос к базе данных для получения фильма по ID
	var movie Movie
	err := db.QueryRow("SELECT id, title, year, description, actors, directors, subtitles, image_path FROM movies WHERE id = $1", id).
		Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Description, &movie.Actors, &movie.Directors, &movie.Subtitles, &movie.ImagePath)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Фильм не найден"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при получении данных из базы данных"})
		}
		return
	}

	// Возвращаем данные фильма в формате JSON
	c.JSON(http.StatusOK, movie)
}

func handleAddMovie(c *gin.Context) {
	// Получение данных из формы
	title := c.PostForm("movieTitle")
	yearStr := c.PostForm("movieYear")
	description := c.PostForm("movieDescription")
	actors := c.PostForm("movieActors")
	directors := c.PostForm("movieDirectors")
	subtitles := c.PostForm("movieSubtitles")

	// Проверка обязательных полей
	if title == "" || yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Название фильма и год выхода обязательны"})
		return
	}

	// Преобразование года в число
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат года"})
		return
	}

	// Обработка загрузки изображения
	file, err := c.FormFile("image")
	var imagePath string
	if err != nil {
		// Если изображение не загружено, используем изображение по умолчанию
		imagePath = "/uploads/default_image.jpg" // Путь к изображению по умолчанию
	} else {
		// Сохранение изображения на сервере
		imagePath = filepath.Join("uploads", file.Filename)
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сохранения изображения"})
			return
		}
		// Возвращаем URL для загрузки изображения
		imagePath = "/uploads/" + file.Filename
	}

	// Вставка данных в базу данных
	query := `
		INSERT INTO movies (title, year, description, actors, directors, subtitles, image_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = db.Exec(query, title, year, description, actors, directors, subtitles, imagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка сохранения данных в базе данных"})
		return
	}

	// Возвращаем путь к изображению (по умолчанию или загруженному)
	c.JSON(http.StatusOK, gin.H{"message": "Фильм успешно добавлен", "image_path": imagePath})
}

func handleGetMovies(c *gin.Context) {
	// Запрос к базе данных для получения всех фильмов
	rows, err := db.Query("SELECT id, title, year, description, actors, directors, subtitles, image_path FROM movies")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при получении данных из базы данных"})
		return
	}
	defer rows.Close()

	var movies []Movie

	// Обработка результатов запроса
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Description, &movie.Actors, &movie.Directors, &movie.Subtitles, &movie.ImagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при обработке данных"})
			return
		}
		movies = append(movies, movie)
	}

	// Возвращаем данные в формате JSON
	c.JSON(http.StatusOK, movies)
}
