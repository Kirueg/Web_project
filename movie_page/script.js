const email = localStorage.getItem('email');
const token = localStorage.getItem('token');
const isLogg = localStorage.getItem('isLoggedIn');
const login = localStorage.getItem('login');
const userId = localStorage.getItem('userId');
console.log(isLogg);
console.log(email);
console.log(token);
console.log(login);
console.log(userId);

if (isLogg === 'true') {
    const lks = document.getElementById('enter_lks');
    lks.textContent = 'Профиль';
    lks.href = '../edit_profile/profile.html';
    const ava = document.getElementsByClassName('avatarka')[0];
}

const urlParams = new URLSearchParams(window.location.search);
const movieID = urlParams.get("id");
if (!movieID) {
    throw new Error("Movie ID is missing");
}

document.addEventListener("DOMContentLoaded", async () => {
    try {
        // Получаем ID фильма из URL
        const urlParams = new URLSearchParams(window.location.search);
        const movieID = urlParams.get("id");
        if (!movieID) {
            throw new Error("Movie ID is missing");
        }

        // Запрос данных с сервера
        const response = await fetch(`http://localhost:8081/api/movies/${movieID}`);
        if (!response.ok) {
            throw new Error(`Ошибка HTTP: ${response.status}`);
        }

        const movie = await response.json();
        console.log("Данные с сервера:", movie); // Вывод данных в консоль

        // Заполняем данные на странице
        document.title = movie.title; // Устанавливаем заголовок страницы
        document.querySelector(".content__image").src = "http://localhost:8081" + movie.image_path;
        document.querySelector(".content__text h2").textContent = `${movie.title} (${movie.year})`;
        document.querySelector(".content__text p").textContent = movie.description;

        // Заполняем актеров
        const actorsList = document.querySelector(".mainActor .Actors");
        actorsList.innerHTML = ""; // Очищаем список
        movie.actors.split(",").forEach((actor) => {
            const li = document.createElement("li");
            li.className = "Actor";
            li.innerHTML = `<a class="Actor_link" href="#">${actor.trim()}</a>`;
            actorsList.appendChild(li);
        });

        // Заполняем режиссеров
        const directorsList = document.querySelector(".director .Actors");
        directorsList.innerHTML = ""; // Очищаем список
        movie.directors.split(",").forEach((director) => {
            const li = document.createElement("li");
            li.className = "Actor";
            li.innerHTML = `<a class="Actor_link" href="#">${director.trim()}</a>`;
            directorsList.appendChild(li);
        });

        // Заполняем субтитры
        document.querySelector(".Sub p").textContent = movie.subtitles;

        // Привязка события на кнопку
        document.getElementById("submitReview").addEventListener("click", () => {
            if (isLogg != 'true') {
                alert('Войдите, чтобы оставлять рецензии!');
            } else {
                submitReview(movieID); // Передаем movieID из URL
            }
        });

        // Отображаем отзывы для текущего фильма
        displayReviews(movieID);
    } catch (error) {
        console.error("Ошибка:", error);
    }
});

async function submitReview(movieID) {
    const reviewText = document.getElementById("reviewText").value;
    const ratingInput = document.querySelector('input[name="hsr"]:checked');

    if (!ratingInput) {
        alert("Пожалуйста, выберите рейтинг.");
        return;
    }

    const rating = parseFloat(ratingInput.value);

    // Проверка, что movieID является числом
    if (isNaN(movieID)) {
        alert("Некорректный ID фильма.");
        return;
    }

    try {
        // Проверка, существует ли уже отзыв от текущего пользователя для этого фильма
        const existingReviewResponse = await fetch(`http://localhost:8081/api/reviews?movie_id=${movieID}&user_id=${userId}`, {
            headers: {
                "Authorization": `Bearer ${token}` // Передаем токен
            }
        });
        if (!existingReviewResponse.ok) {
            const errorData = await existingReviewResponse.json();
            throw new Error(`Ошибка HTTP: ${existingReviewResponse.status}, ${errorData.message || "Неизвестная ошибка"}`);
        }

        const existingReviews = await existingReviewResponse.json();
        if (existingReviews && Array.isArray(existingReviews)) {
            if (existingReviews.length > 0) {
                // Если отзыв уже существует, обновляем его
                updateReview(existingReviews[0].id, reviewText, rating);
            } else {
                // Если отзыва нет, добавляем новый
                addReview(movieID, reviewText, rating);
            }
        } else {
            addReview(movieID, reviewText, rating);
        }
    } catch (error) {
        addReview(movieID, reviewText, rating);
    }
}

function addReview(movieID, reviewText, rating) {
    fetch("http://localhost:8081/api/reviews", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}` // Передаем токен
        },
        body: JSON.stringify({
            author: login, // Убедитесь, что имя автора передается
            text: reviewText,
            rating: rating,
            movie_id: parseInt(movieID), // Убедитесь, что movieID — целое число
            user_id: parseInt(userId), // Передаем user_id
        }),
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(errorData => {
                    const errorMessage = `Server error: ${errorData.message || JSON.stringify(errorData)}`;
                    throw new Error(errorMessage);
                }).catch(err => {
                    const errorMessage = `HTTP error! status: ${response.status}, details: ${err.message}`;
                    throw new Error(errorMessage);
                });
            }
            return response.json();
        })
        .then(data => {
            console.log("Review added:", data);
            displayReviews(movieID);
            document.getElementById("reviewText").value = "";
        })
        .catch(error => {
            console.error("Error adding review:", error);
            alert(`Error adding review: ${error.message}`); // Show the actual error
        });
}

function updateReview(reviewID, reviewText, rating) {
    fetch(`http://localhost:8081/api/reviews/${reviewID}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}` // Передаем токен
        },
        body: JSON.stringify({
            text: reviewText,
            rating: rating,
        }),
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(errorData => {
                    const errorMessage = `Server error: ${errorData.message || JSON.stringify(errorData)}`;
                    throw new Error(errorMessage);
                }).catch(err => {
                    const errorMessage = `HTTP error! status: ${response.status}, details: ${err.message}`;
                    throw new Error(errorMessage);
                });
            }
            return response.json();
        })
        .then(data => {
            console.log("Review updated:", data);
            displayReviews(movieID);
            document.getElementById("reviewText").value = "";
        })
        .catch(error => {
            console.error("Error updating review:", error);
            alert(`Error updating review: ${error.message}`); // Show the actual error
        });
}

document.addEventListener("DOMContentLoaded", () => {
    // Обработчик для удаления комментариев
    document.getElementById("reviewsContainer").addEventListener("click", async (event) => {
        if (event.target.closest(".del_rev")) {
            const reviewElement = event.target.closest(".del_rev").parentElement;
            const reviewUserId = parseInt(reviewElement.dataset.userId, 10); // ID автора комментария
            const reviewID = reviewElement.dataset.reviewId; // ID комментария

            // Проверяем, соответствует ли автор комментария текущему пользователю
            if (reviewUserId !== parseInt(userId, 10)) {
                alert("Вы можете удалять только свои комментарии.", reviewUserId, userId, 2);
                return;
            }

            try {
                // Удаляем комментарий с сервера
                const response = await fetch(`http://localhost:8081/api/reviews/${reviewID}`, {
                    method: "DELETE",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${token}` // Передаем токен
                    },
                });

                if (!response.ok) {
                    throw new Error(`Ошибка HTTP: ${response.status}`);
                }

                // Удаляем комментарий из DOM
                reviewElement.remove();
            } catch (error) {
                console.error("Ошибка при удалении комментария:", error);
                alert("Ошибка при удалении комментария: " + error.message);
            }
        }
    });
});

function displayReviews(movieID) {
    document.getElementById("reviewsContainer").innerHTML = ""; // Clear existing reviews

    fetch(`http://localhost:8081/api/reviews?movie_id=${movieID}`, {
        headers: {
            "Authorization": `Bearer ${token}` // Передаем токен
        }
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(errorData => {
                    const errorMessage = `Server error: ${errorData.message || JSON.stringify(errorData)}`;
                    throw new Error(errorMessage);
                }).catch(err => {
                    const errorMessage = `HTTP error! status: ${response.status}, details: ${err.message}`;
                    throw new Error(errorMessage);
                });
            }
            return response.json();
        })
        .then(reviews => {
            // Проверка, что reviews является массивом
            if (!Array.isArray(reviews)) {
                return;
            }

            // Если отзывов нет, выводим сообщение
            if (reviews.length === 0) {
                const noReviewsElement = document.createElement("div");
                noReviewsElement.textContent = "Отзывов пока нет.";
                document.getElementById("reviewsContainer").appendChild(noReviewsElement);
                return;
            }

            reviews.forEach(review => {
                const reviewElement = document.createElement("div");
                reviewElement.dataset.reviewId = review.id; // ID комментария
                reviewElement.dataset.userId = review.user_id; // ID автора комментария

                reviewElement.innerHTML = `
                    <span class="author_name">${review.author}:</span>
                    <button class="del_rev">
                        <span class="xcross">&times;</span>
                    </button>
                    <p>${review.text}</p>
                    <span>Оценка: ${review.rating}</span>
                `;
                document.getElementById("reviewsContainer").appendChild(reviewElement);
            });
        })
        .catch(error => {
            console.error("Error fetching reviews:", error);
            alert(`Error fetching reviews: ${error.message}`); // Show the actual error message
        });
}