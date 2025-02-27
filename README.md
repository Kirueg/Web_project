# Кинокритика - Веб-сайт для написания рецензий на фильмы

Добро пожаловать в репозиторий веб-сайта для написания рецензий на фильмы! Этот проект предоставляет пользователям возможность делиться своими мнениями о фильмах, оценивать их и читать отзывы других пользователей. Серверная часть проекта реализована на языке программирования Golang с использованием микросервисной архитектуры.

## Главная страница сайта

Ниже представлен скриншот главной страницы сайта:

<img width="914" alt="image" src="https://github.com/user-attachments/assets/6c7d9571-f2d4-4b48-8727-cdd04be63a64" />

## Основные функции

Проект состоит из нескольких микросервисов, каждый из которых отвечает за определённую функциональность:

1. **Управление карточками фильмов:**
   - Добавление новых карточек фильмов.
   - Удаление существующих карточек фильмов.

2. **Управление отзывами и оценками:**
   - Добавление отзывов и оценок к фильмам.
   - Редактирование существующих отзывов и оценок.
   - Удаление отзывов и оценок.

3. **Авторизация пользователей:**
   - Регистрация новых пользователей.
   - Аутентификация пользователей.
   - Управление сессиями пользователей.

## Технологии

- **Язык программирования:** Golang (для серверной части).
- **База данных:** PostgreSQL (для хранения данных о фильмах, пользователях, отзывах и оценках).
- **Микросервисы:** Архитектура проекта построена на основе микросервисов, что обеспечивает масштабируемость и гибкость.
