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

if(isLogg === 'true'){
    const lks = document.getElementById('enter_lks');
    lks.textContent = 'Профиль';
    lks.href = '../edit_profile/profile.html';
    const ava = document.getElementsByClassName('avatarka')[0];
    ava.style.display = 'inline-block';
}

document.addEventListener("DOMContentLoaded", async () => {
    try {
        const container = document.querySelector(".card-container");
        const searchInput = document.getElementById("Search");

        // Функция для загрузки фильмов
        const loadMovies = async (searchTerm = "") => {
            const response = await fetch(`http://localhost:8081/api/movies?search=${searchTerm}`);
            if (!response.ok) {
                throw new Error("Ошибка при загрузке данных с сервера");
            }

            const movies = await response.json();

            // Очищаем контейнер перед добавлением новых карточек
            container.innerHTML = "";

            // Создаем карточки для каждого фильма
            movies.forEach((movie) => {
                const card = document.createElement("div");
                card.classList.add("card");

                const link = document.createElement("a");
                link.href = `../movie_page/movie_template.html?id=${movie.id}`; // Переход на страницу фильма
                link.classList.add("image-card");

                const image = document.createElement("img");
                image.src = "http://localhost:8081" + movie.image_path; // Путь к изображению
                image.alt = movie.title;
                link.appendChild(image);

                // Добавляем крестик для удаления
                const deleteIcon = document.createElement("div");
                deleteIcon.classList.add("delete-icon");
                deleteIcon.innerHTML = "<span>&times;</span>";
                deleteIcon.dataset.id = movie.id; // Сохраняем id фильма в data-атрибут

                // Обработчик события для удаления карточки
                deleteIcon.addEventListener("click", async (event) => {
                    event.preventDefault(); // Предотвращаем переход по ссылке
                    const movieID = deleteIcon.dataset.id;

                    // Отправляем запрос на сервер для удаления фильма
                    const deleteResponse = await fetch(`http://localhost:8081/api/movies/${movieID}`, {
                        method: "DELETE",
                    });

                    if (deleteResponse.ok) {
                        // Удаляем карточку из DOM
                        card.remove();
                        console.log("Фильм удален:", movieID);
                    } else {
                        console.error("Ошибка при удалении фильма:", deleteResponse.status);
                    }
                });

                card.appendChild(link);
                card.appendChild(deleteIcon); // Добавляем крестик
                container.appendChild(card);
            });
        };

        // Загружаем все фильмы при загрузке страницы
        loadMovies();

        // Обработчик события для поиска
        searchInput.addEventListener("input", () => {
            const searchTerm = searchInput.value.trim();
            loadMovies(searchTerm);
        });

    } catch (error) {
        console.error("Ошибка:", error);
    }
});



