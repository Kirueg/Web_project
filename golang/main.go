package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

type Review struct {
	ID      int     `json:"id"`
	Text    string  `json:"text"`
	Author  string  `json:"author"`
	Rating  float64 `json:"rating"`
	MovieID int     `json:"movie_id"`
	UserID  int     `json:"user_id"` // Добавлено поле user_id
}

type Account struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

var db *sql.DB

const (
	host      = "localhost"
	port      = 5432
	user      = "said"
	password  = "12102005"
	dbname    = "kinopoisk"
	jwtSecret = "your_secret_key" // Используйте более сложный секрет
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

	// Создание таблиц, если они не существуют
	createTables()

	// Настройка маршрутов
	r := gin.Default()

	// Подключаем CORS-мидлвар
	r.Use(corsMiddleware())

	// Маршрут для статических файлов (изображений)
	r.Static("/uploads", "./uploads")

	// Маршруты для фильмов
	r.GET("/api/movies/:id", getMovieByID)
	r.POST("/api/add-movie", handleAddMovie)
	r.GET("/api/movies", handleGetMovies)
	r.DELETE("/api/movies/:id", deleteMovie)

	// Маршруты для отзывов
	r.GET("/api/reviews", getReviews)
	r.POST("/api/reviews", addReview) //authMiddleware(),
	r.PUT("/api/reviews/:id", updateReview)
	r.DELETE("/api/reviews/:id", deleteReview)

	// Маршруты для аккаунтов
	r.GET("/aaccounts", getAccounts)
	r.POST("/account", postAccount)
	r.POST("/login", postLogin)
	r.GET("/login-by-email", getLoginByEmail)
	r.POST("/update-profile", updateProfile)

	// Запуск сервера
	fmt.Println("Сервер запущен на http://localhost:8081")
	r.Run(":8081")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func deleteMovie(c *gin.Context) {
	movieID := c.Param("id")

	// Удаление связанных записей из таблицы reviews
	_, err := db.Exec("DELETE FROM reviews WHERE movie_id = $1", movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при удалении отзывов", "error": err.Error()})
		return
	}

	// Удаление фильма из таблицы movies
	_, err = db.Exec("DELETE FROM movies WHERE id = $1", movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при удалении фильма", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Фильм успешно удален", "id": movieID})
}

func deleteReview(c *gin.Context) {
	reviewID := c.Param("id")

	// Удаление комментария из базы данных
	_, err := db.Exec("DELETE FROM reviews WHERE id = $1", reviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при удалении комментария", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Комментарий успешно удален", "id": reviewID})
}

func createTables() {
	// Создание таблицы для фильмов
	queryMovies := `
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
	_, err := db.Exec(queryMovies)
	if err != nil {
		fmt.Println("Ошибка создания таблицы movies:", err)
	}

	// Создание таблицы для отзывов
	queryReviews := `
		CREATE TABLE IF NOT EXISTS reviews (
			id SERIAL PRIMARY KEY,
			text TEXT NOT NULL,
			author TEXT NOT NULL,
			rating FLOAT NOT NULL,
			movie_id INT NOT NULL,
			user_id INT NOT NULL
		)
	`
	_, err = db.Exec(queryReviews)
	if err != nil {
		fmt.Println("Ошибка создания таблицы reviews:", err)
	}

	// Создание таблицы для аккаунтов
	queryAccounts := `
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			login TEXT NOT NULL,
			password TEXT NOT NULL,
			email TEXT NOT NULL
		)
	`
	_, err = db.Exec(queryAccounts)
	if err != nil {
		fmt.Println("Ошибка создания таблицы accounts:", err)
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
	searchTerm := c.Query("search")

	var movies []Movie
	var query string
	var args []interface{}

	if searchTerm != "" {
		query = `SELECT id, title, image_path FROM movies WHERE title ILIKE $1`
		args = []interface{}{"%" + searchTerm + "%"}
	} else {
		query = `SELECT id, title, image_path FROM movies`
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при выполнении запроса"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.ImagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сканировании данных"})
			return
		}
		movies = append(movies, movie)
	}

	c.JSON(http.StatusOK, movies)
}

func addReview(c *gin.Context) {
	var newReview Review
	err := c.BindJSON(&newReview)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат данных", "error": err.Error()})
		return
	}

	// Получаем userID из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Пользователь не авторизован"})
		return
	}
	newReview.UserID = userID.(int)

	// Проверка, что movie_id существует в таблице movies
	var movieExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM movies WHERE id = $1)", newReview.MovieID).Scan(&movieExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка проверки movie_id", "error": err.Error()})
		return
	}
	if !movieExists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Указанный movie_id не существует"})
		return
	}

	// Вставка данных в базу данных
	stmt, err := db.Prepare("INSERT INTO reviews (text, author, rating, movie_id, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при сохранении данных", "error": err.Error()})
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(newReview.Text, newReview.Author, newReview.Rating, newReview.MovieID, newReview.UserID).Scan(&newReview.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при сохранении данных", "error": err.Error()})
		return
	}

	// Возвращаем данные в формате JSON
	c.JSON(http.StatusOK, newReview)
}

func updateReview(c *gin.Context) {
	reviewID := c.Param("id")
	var updatedReview Review
	err := c.BindJSON(&updatedReview)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат данных", "error": err.Error()})
		return
	}

	// Обновление данных в базе данных
	stmt, err := db.Prepare("UPDATE reviews SET text = $1, rating = $2 WHERE id = $3")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при обновлении данных", "error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedReview.Text, updatedReview.Rating, reviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при обновлении данных", "error": err.Error()})
		return
	}

	// Возвращаем данные в формате JSON
	c.JSON(http.StatusOK, gin.H{"message": "Отзыв успешно обновлен", "id": reviewID})
}

func getReviews(c *gin.Context) {
	movieID := c.Query("movie_id")
	userID := c.Query("user_id")

	if movieID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Не указан movie_id"})
		return
	}

	var rows *sql.Rows
	var err error

	if userID != "" {
		// Если указан user_id, фильтруем по user_id и movie_id
		rows, err = db.Query("SELECT id, text, author, rating, user_id FROM reviews WHERE movie_id = $1 AND user_id = $2", movieID, userID)
	} else {
		// Иначе получаем все отзывы для фильма
		rows, err = db.Query("SELECT id, text, author, rating, user_id FROM reviews WHERE movie_id = $1", movieID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при получении данных из базы данных", "error": err.Error()})
		return
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		err = rows.Scan(&review.ID, &review.Text, &review.Author, &review.Rating, &review.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при обработке данных", "error": err.Error()})
			return
		}
		reviews = append(reviews, review)
	}

	// Возвращаем данные в формате JSON
	c.JSON(http.StatusOK, reviews) // Всегда возвращаем массив
}

func getAccounts(c *gin.Context) {
	rows, err := db.Query("SELECT id, login, password FROM accounts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при получении данных из базы данных", "error": err.Error()})
		return
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var acc Account
		if err := rows.Scan(&acc.ID, &acc.Login, &acc.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при обработке данных", "error": err.Error()})
			return
		}
		accounts = append(accounts, acc)
	}

	c.JSON(http.StatusOK, accounts)
}

func postAccount(c *gin.Context) {
	var input Account
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат данных", "error": err.Error()})
		return
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка хэширования пароля", "error": err.Error()})
		return
	}
	input.Password = string(hashedPassword)

	// Вставка данных в базу данных
	_, err = db.Exec("INSERT INTO accounts (login, password, email) VALUES ($1, $2, $3)", input.Login, input.Password, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при сохранении данных", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Аккаунт успешно создан"})
}

func postLogin(c *gin.Context) {
	var input Account
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат данных", "error": err.Error()})
		return
	}

	var storedPassword string
	var userId int
	err = db.QueryRow("SELECT id, password FROM accounts WHERE email = $1", input.Login).Scan(&userId, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Логин или пароль неверны"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при получении данных из базы данных", "error": err.Error()})
		return
	}

	// Сравнение введенного пароля с захэшированным
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Логин или пароль неверны"})
		return
	}

	// Создаем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": input.Login,
		"id":    userId, // Добавляем ID пользователя в токен
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при создании токена", "error": err.Error()})
		return
	}

	// Отправляем токен
	response := map[string]string{
		"token": tokenString,
		"login": input.Login,
	}

	c.JSON(http.StatusOK, response)
}

func getLoginByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Отсутствует параметр email"})
		return
	}

	var account Account
	err := db.QueryRow("SELECT id, login FROM accounts WHERE email = $1", email).Scan(&account.ID, &account.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Аккаунт с таким email не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при получении данных из базы данных", "error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"id":    account.ID,
		"login": account.Login,
	}

	c.JSON(http.StatusOK, response)
}

func updateProfile(c *gin.Context) {
	var input Account
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат данных", "error": err.Error()})
		return
	}

	// Проверяем, что ID пользователя передан
	if input.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Отсутствует ID пользователя"})
		return
	}

	// Хэширование нового пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка хэширования пароля", "error": err.Error()})
		return
	}
	input.Password = string(hashedPassword)

	// Обновляем данные пользователя по ID
	_, err = db.Exec("UPDATE accounts SET email = $1, login = $2 WHERE id = $3", input.Email, input.Login, input.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при обновлении профиля", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Профиль успешно обновлен"})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		// fmt.Println(authHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Отсутствует токен авторизации"})
			return
		}

		// Извлекаем токен из заголовка
		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем метод подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Недействительный токен", "error": err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Добавляем информацию о пользователе в контекст
			userID := int(claims["id"].(float64))
			c.Set("userID", userID)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Недействительный токен"})
			return
		}

		c.Next()
	}
}

// http://localhost:8081/api/reviews?movie_id=7&user_id=111111&Text=5252525252552
// http://localhost:8081/login?Login=aaa@aaa.aaa
