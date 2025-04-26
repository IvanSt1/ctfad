package routing

import (
    "os"
    "strings"

    "github.com/gorilla/schema"
    "github.com/gorilla/sessions"

    // Правильный путь вашего модуля
    "github.com/IvanSt1/ctfad/otkritki/backend/core/db"
)

const (
    cookieName  = "session"
    authKey     = "authenticated"
    genderKey   = "gender"
    idKey       = "id"

    usernameLen = 10
    passwordLen = 10
)

var (
    // Подключение к БД (предполагается функция GetDB в core/db)
    database = db.GetDB()

    // CookieStore с ключами из .env
    cookieStore = sessions.NewCookieStore(
        []byte(os.Getenv("COOKIE_AUTH_KEY")),
        []byte(os.Getenv("COOKIE_ENC_KEY")),
    )

    // Схема encoder/decoder для разбора форм
    encoder = schema.NewEncoder()
    decoder = schema.NewDecoder()

    // Белый список Origin для CORS
    allowedOrigins = strings.Split(os.Getenv("CORS_ALLOWED"), ",")
)
