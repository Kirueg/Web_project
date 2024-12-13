document.addEventListener("DOMContentLoaded", async () => {
    try {
        // Получаем ID фильма из URL
        const urlParams = new URLSearchParams(window.location.search);
        const movieID = urlParams.get("id");
        if (!movieID) {
            throw new Error("Movie ID is missing");
        }

        // Запрос данных с сервера
        const response = await fetch(`http://localhost:8080/api/movies/${movieID}`);
        if (!response.ok) {
            throw new Error(`Ошибка HTTP: ${response.status}`);
        }

        const movie = await response.json();
        console.log("Данные с сервера:", movie); // Вывод данных в консоль

        // Заполняем данные на странице
        document.title = movie.title; // Устанавливаем заголовок страницы
        document.querySelector(".content__image").src = "http://localhost:8080" + movie.image_path;
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
    } catch (error) {
        console.error("Ошибка:", error);
    }
});