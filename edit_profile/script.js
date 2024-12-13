// Функция для управления состоянием
function useState(initialValue) {
    let state = initialValue;

    function setState(newValue) {
        state = newValue;
    }

    return [state, setState];
}

// Основная функция
function App() {
    const [isLoggedIn, setIsLoggedIn] = useState(false); // Инициализация состояния

    localStorage.removeItem('token');
    localStorage.removeItem('email');
    localStorage.removeItem('isLoggedIn');
    localStorage.removeItem('login');
    setIsLoggedIn(false); // Обновляем состояние
    window.location.href = '../main_page/index.html'; // Перенаправление

    return;
}

const logg = localStorage.getItem('login');
const email = localStorage.getItem('email');
const nickname = document.getElementById('nickname')
const emaill = document.getElementById('email')
nickname.textContent = 'Логин: ' + logg;
emaill.textContent = 'Email: ' + email;


emailInput.value = email;
nameInput.value = logg;

