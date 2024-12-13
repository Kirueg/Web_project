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
        // Запрос данных с сервера
        const response = await fetch("http://localhost:8080/api/movies");
        if (!response.ok) {
            throw new Error("Ошибка при загрузке данных с сервера");
        }

        const movies = await response.json();

        // Контейнер для карточек фильмов
        const container = document.querySelector(".card-container");

        // Создаем карточки для каждого фильма
        movies.forEach((movie) => {
            const card = document.createElement("div");
            card.classList.add("card");

            const link = document.createElement("a");
            link.href = `../movie_page/movie_template.html?id=${movie.id}`; // Переход на страницу фильма
            link.classList.add("image-card");

            const image = document.createElement("img");
            image.src = "http://localhost:8080" + movie.image_path; // Путь к изображению
            image.alt = movie.title;
            link.appendChild(image);

            card.appendChild(link);
            container.appendChild(card);
        });
    } catch (error) {
        console.error("Ошибка:", error);
    }
});