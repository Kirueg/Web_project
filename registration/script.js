// import React, { useState } from 'react';
// import ReactDOM from 'react-dom';
// export default LoginForm;

// const LoginForm = ({ onTokenReceived, currentLogin }) => {
//   const [login, setLogin] = useState('');
//   const [password, setPassword] = useState('');
//   const [errorMessage, setErrorMessage] = useState('');
//   const [isFormVisible, setFormVisible] = useState(false);
//   const [isLoggedIn, setLoggedIn] = useState(false);

//   const handleLoginClick = () => {
//     setFormVisible(true);
//   };

//   const handleCancelClick = () => {
//     setFormVisible(false);
//     setErrorMessage('');
//     setLogin('');
//     setPassword('');
//   };

//   const handleSubmit = async (e) => {
//     e.preventDefault();
//     setErrorMessage('');

//     const response = await fetch('http://localhost:8081/login', {
//       method: 'POST',
//       headers: {
//         'Content-Type': 'application/json',
//       },
//       body: JSON.stringify({ login, password }),
//     });

//     if (response.ok) {
//       const data = await response.json();
//       alert('Авторизация успешна! Токен: ' + data.token);
//       setLoggedIn(true);
//       setFormVisible(false);
//       onTokenReceived(data.token, data.login, data.password); // Передаем логин и пароль в родительский компонент
//     } else {
//       const errorText = await response.text();
//       setErrorMessage('Ошибка: ' + errorText);
//     }
//   };

//   return (
//     <div>
//       {isLoggedIn ? (
//         <p>Вы авторизованы как: {currentLogin}</p> // Отображаем текущий логин
//       ) : (
//         <>
//           {!isFormVisible ? (
//             <button onClick={handleLoginClick}>Вход</button> // Кнопка для входа
//           ) : (
//             <div>
//               <h1>Авторизация</h1>
//               <form onSubmit={handleSubmit}>
//                 <label htmlFor="login">Логин:</label>
//                 <input
//                   type="text"
//                   id="login"
//                   value={login}
//                   onChange={(e) => setLogin(e.target.value)}
//                   required
//                 />
//                 <br />
//                 <label htmlFor="password">Пароль:</label>
//                 <input
//                   type="password"
//                   id="password"
//                   value={password}
//                   onChange={(e) => setPassword(e.target.value)}
//                   required
//                 />
//                 <br />
//                 <button type="submit">Войти</button>
//                 <button type="button" onClick={handleCancelClick}>Отмена</button> {/* Кнопка Отмена */}
//               </form>
//               {errorMessage && <p style={{ color: 'red' }}>{errorMessage}</p>}
//             </div>
//           )}
//         </>
//       )}
//     </div>
//   );
// };

// Используем глобальную переменную
// Читаем данные из localStorage
const myVariable = localStorage.getItem('myVariable');
console.log(myVariable); // Выведет: Hello from index1.html!

document.querySelector('.add_ava').addEventListener('change', function(event) {
    const file = event.target.files[0]; // Получаем выбранный файл

    if (file && file.type.startsWith('image/')) { // Проверяем, что это изображение
        const reader = new FileReader(); // Создаем объект FileReader

        reader.onload = function(e) {
            const previewImage = document.getElementById('preview_ava');
            previewImage.src = e.target.result; // Устанавливаем src для изображения
            previewImage.style.display = 'block'; // Показываем изображение
        };

        reader.readAsDataURL(file); // Читаем файл как Data URL
    } else {
        alert('Пожалуйста, выберите изображение.');
    }
});

function goBack() {
    window.history.back(); // Возвращает на предыдущую страницу
}

function goBackx2() {
    window.history.go(-2);
}


document.getElementById('registrationForm').addEventListener('submit', function (event) {
  event.preventDefault(); // Предотвращаем отправку формы по умолчанию

  const login = document.getElementById('login').value.trim();
  const email = document.getElementById('email').value.trim();
  const password = document.getElementById('password').value.trim();
  const confirmPassword = document.getElementById('confirmPassword').value.trim();

  if (!login || !email || !password || !confirmPassword) {
    document.getElementById('errorMessage').textContent = 'Пожалуйста, заполните все поля';
    return; // Прерываем выполнение, если есть пустые поля
  }

  // Проверка на совпадение паролей
  if (password !== confirmPassword) {
    document.getElementById('errorMessage').textContent = 'Пароли не совпадают';
    return;
  }

  // Отправка данных на сервер
  fetch('http://localhost:8081/account', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ login, email, password }),
  })
    .then(response => {
      if (response.ok) {
        // После успешной регистрации получаем логин и ID по email
        return fetch(`http://localhost:8081/login-by-email?email=${email}`)
          .then(response => {
            if (response.ok) {
              return response.json(); // Преобразуем ответ в JSON
            } else {
              throw new Error('Ошибка при получении логина и ID по email');
            }
          })
          .then(data => {
            // Сохраняем логин, email и ID в localStorage
            localStorage.setItem('login', data.login);
            localStorage.setItem('email', email);
            localStorage.setItem('userId', data.id); // Сохраняем ID
            localStorage.setItem('isLoggedIn', 'true'); // Устанавливаем флаг входа
            window.location.href = '../main_page/index.html'; // Перенаправление на главную страницу
          });
      } else {
        return response.text().then(errorText => {
          document.getElementById('errorMessage').textContent = 'Ошибка: ' + errorText;
        });
      }
    })
    .catch(error => {
      document.getElementById('errorMessage').textContent = 'Ошибка сети: ' + error.message;
    });
});