# go-auth-service
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com) [![forthebadge](http://forthebadge.com/images/badges/built-with-love.svg)](http://forthebadge.com)

Сервис для ведения БД пользователей и контроля их сессий.
Используются принципы чистой архитектуры.

Доступ к сервису осуществляется только через запросы подписанные **JWT токеном** которые генерируется при запуске и сохраняется в ENV.

# Запуск
Скачать репозиторий:
```
git clone https://github.com/giicoo/go-auth-service
cd go-auth-service
```

Собрать образ докера:
```
docker build -t auth:1.0 .
```

Запустить docker-compose:
```
docker-compose up --build
```
Swagger-docs:
http://localhost:8080/swagger/index.html

# ***Endpoints:***
## Users
#### /create-user
##### POST
*request:*
```json
{
	"email": "string",
	"password": "string"
}
```

*response:*
```json
{
	"id": 0,
	"email": "string",
}
```

##### **Engine**
Проверка отсутствия аккаунта с этим email. Пароль хешируется с помощью `bcrypt`. *User* записывается записывается в БД  

#### /delete-user
##### DELETE
*request:*
```json
{
    "id":0
}
```

*response:*
```json
{
	"message": "user successfully deleted"
}
```

##### **Engine**
Удаление из БД по ID из Session

#### /update-email
##### PUT
*request:*

```json
{
    	"id":0,
	"new_email": "string"
}
```

*response:*
```json
{
	"message": "user success update"
}
```

##### Engine
Проверка не занят ли новый email. Обновление полей в БД по ID из Session

#### /update-password
##### PUT
*request:*
```json
{
    	"id":0,
	"new_password": "string"
}
```

*response:*
```json
{
	"message": "user success update"
}
```

##### Engine
Хеширование нового пароля с помощью bcrypt. Обновление полей в БД по ID из Session

#### /get-user-by-id/{id}
##### GET
*request:*
```
id: integer
```

*response:*
```json
{
	"id": 0,
	"email": "string",
}
```

##### Engine
Получение user из БД

#### /check-user
##### POST
*request:*
```json
{
	"email": "string",
	"password": "string"
}
```

*response:*
```json
{
	"id": 0,
	"email": "string",
}
```

##### **Engine**
Проверка пароля


---

## Sessions
**Redis**

#### /create-session
##### POST
*request:*
```json
{
	"user_id": 0,
	"user_ip": "string",
	"user_agent": "string"
}
```

*response:*
```json
{
	"session_id": "string",
}
```

##### **Engine**
SessionID генерируется как 32 рандомных бит
Добавляется в Redis:
1) *session <session_id>{"user_id": int, "user_ip": "string", "user_agent": "string"}*
2) *user_sessions <user_id>["session_id",...]*

#### /get-session/{session_id}
##### GET
*request:*

```
session_id: string
```

*response:*
```json
{
	"session_id": "string",
	"user_id": 0,
	"user_ip": "string",
	"user_agent": "string"
}
```

##### Engine
Получение из Redis: session <session_id>



##### /get-sessions/{user_id}
##### GET
*request:*

```
user_id: int
```

*response:*
```json

[
	{
		"session_id": "string",
		"user_id": 0,
		"user_ip": "string",
		"user_agent": "string"
	},
	...
]

```
##### Engine
Получение из Redis: user_sessions <user_id>


##### /delete-session
##### DELETE
*request:*

```json
{
	"session_id": "string"
}
```

*response:*
```json
{
	"message": "session successfully deleted"
}
```

##### Engine
Удаление сессии session_id из Redis: session <session_id> если у сессий совпадает user_id
Удаление сессии из списка сессий юзера

##### /delete-sessions
##### DELETE
*request:*

```json
{
	"user_id": "string"
}
```

*response:*
```json
{
	"message": "sessions successfully deleted"
}
```

##### Engine
Удаление сессий по user_id из Redis: user_sessions <user_id>
