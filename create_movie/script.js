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
        }


// Обработчик изменения файла (загрузка изображения)
document.getElementById('fileInput').addEventListener('change', function(event) {
    const file = event.target.files[0];
    if (file) {
        const reader = new FileReader();
        reader.onload = function(e) {
            // Создаем объект изображения
            const img = new Image();
            img.src = e.target.result;

            // Обработчик загрузки изображения
            img.onload = function() {
                // Создаем canvas для масштабирования изображения
                const canvas = document.createElement('canvas');
                const ctx = canvas.getContext('2d');

                // Устанавливаем размеры canvas
                canvas.width = 300; // Фиксированная ширина
                canvas.height = 400; // Фиксированная высота

                // Масштабируем изображение
                ctx.drawImage(img, 0, 0, canvas.width, canvas.height);

                // Обновляем изображение на странице
                document.getElementById('preview_ava').src = canvas.toDataURL('image/jpeg');
            }
        }
        reader.readAsDataURL(file);
    }
});

// Обработчик нажатия на кнопку "Добавить"
document.getElementById('submitButton').addEventListener('click', async () => {
    const movieTitle = document.getElementById('movieTitle').value;
    const movieYear = document.getElementById('movieYear').value;
    const movieDescription = document.getElementById('movieDescription').value;
    const movieActors = document.getElementById('movieActors').value;
    const movieDirectors = document.getElementById('movieDirectors').value;
    const movieSubtitles = document.getElementById('movieSubtitles').value;
    const fileInput = document.getElementById('fileInput').files[0];

    // Проверка на заполненность обязательных полей
    if (!movieTitle || !movieYear) {
        alert('Пожалуйста, заполните обязательные поля: название фильма, год выхода');
        return;
    }

    // Создаем FormData для отправки данных
    const formData = new FormData();
    formData.append('movieTitle', movieTitle);
    formData.append('movieYear', movieYear);
    formData.append('movieDescription', movieDescription);
    formData.append('movieActors', movieActors);
    formData.append('movieDirectors', movieDirectors);
    formData.append('movieSubtitles', movieSubtitles);
    formData.append('image', fileInput);

    try {
        const response = await fetch('http://localhost:8081/api/add-movie', {
            method: 'POST',
            body: formData,
        });

        if (response.ok) {
            const data = await response.json();
            alert('Фильм успешно добавлен!');
            // Очищаем форму
            document.getElementById('movieTitle').value = '';
            document.getElementById('movieYear').value = '';
            document.getElementById('movieDescription').value = '';
            document.getElementById('movieActors').value = '';
            document.getElementById('movieDirectors').value = '';
            document.getElementById('movieSubtitles').value = '';
            document.getElementById('fileInput').value = '';
            // Возвращаем изображение к дефолтному
            document.getElementById('preview_ava').src = "http://localhost:8081/uploads/default_image.jpg";
        } else {
            const errorData = await response.json();
            alert(`Ошибка: ${errorData.message}`);
        }
    } catch (error) {
        console.error('Ошибка при отправке данных:', error);
        alert('Произошла ошибка при отправке данных на сервер.');
    }
});