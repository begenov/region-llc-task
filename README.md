# Region Todo List
Этот микросервис реализует функциональность управления списком задач.

# Требования к запуску
### Предварительные требования
-  Установленный Docker
- Установленный Docker Compose
## Установка и Запуск приложения

1. Клонируйте репозиторий проекта на свой компьютер:
```shell
git clone https://github.com/begenov/region-llc-task
```
2. Перейдите в директорию проекта:
```shell
cd ./region-llc-task
```
3. Запустите приложение с помощью Docker Compose:
```shell
make run
```
или
```shell
docker-compose up
```
4. Чтобы остановить приложение, выполните команду:
```shell
make stop
```
или
```
docker-compose down
```


### API Endpoints
#### ** Формат обмена данными JSON.**
#### Swagger документация доступна по адресу http://localhost:8080/swagger/index.html

## Создание новой Пользователя

1. Метод: POST
- URL: /api/v1/users/sign-up
- Тело запроса:

```json
{
   "email":"test@example.com",
   "username": "username",
   "password": "password"
}
```
- Регистрирует нового пользователя
## Вход пользователя

2. Метод: POST
    
- URL: /api/v1/users/sign-in
- Тело запроса:
```json
{
   "email":"test@example.com",
   "password": "password"
}
```
- Вход пользователя.

## Обновление токена аутентификации

3. Метод: POST

- URL: /api/v1/users/auth/refresh
- Тело запроса:

```json
{
   "refresh_token": "your-refresh-token"
}
```

- Обновляет токен аутентификации.

## Создание задачи

4. Метод: POST
- URL: /api/v1/users/todo-list/tasks
- Авторизация: Bearer "ваш-доступ-токен"
- Тело запроса:

```json
{
   "title": "Купить книгу",
   "activeAt": "2023-08-04"
}
```
- Создание новой задачи

## Обновление задачи

5. Метод: PUT
- URL: /api/v1/users/todo-list/tasks/:id
- Авторизация: Bearer "ваш-доступ-токен"
- Тело запроса:
```json
{
   "title": "Купить книгу - Высоконагруженные приложения",
   "activeAt": "2023-08-05"
}
```
- Обновление существующей задачи.

## Удаление задачи

6. Метод: DELETE
- URL: /api/v1/users/todo-list/tasks/:id
- Авторизация: Bearer "ваш-доступ-токен"

   - Удаление задачи.


## Пометить задачу выполненной

7. Метод: PUT
- URL: /api/v1/users/todo-list/tasks/:id/done
- Авторизация: Bearer "ваш-доступ-токен"

  -  Помечает задачу как выполненную.

## Список задач

8. Метод: GET
- URL: /api/v1/users/todo-list/tasks
- Авторизация: Bearer "ваш-доступ-токен"

  -  Получает список задач.


##  Тестирование
Запуск unit тестов
```shell
go test -cover ./...
```
