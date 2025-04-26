package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/csrf"
    "github.com/gorilla/mux"
    "github.com/gorilla/sessions"

    "github.com/IvanSt1/ctfad/otkritki/backend/core/routing"
)

func main() {
    // Настройка хранилища сессий
    store := sessions.NewCookieStore(
        []byte(os.Getenv("COOKIE_AUTH_KEY")),
        []byte(os.Getenv("COOKIE_ENC_KEY")),
    )
    store.Options = &sessions.Options{
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
        Path:     "/",
    }

    r := mux.NewRouter()

    // CORS middleware из routing
    r.Use(routing.CorsMiddleware)

    // CSRF middleware
    csrfMw := csrf.Protect(
        []byte(os.Getenv("CSRF_KEY")),
        csrf.Secure(true),
    )

    // Регистрация маршрутов
    routing.RegisterRoutes(r)
    routing.RegisterPost(r, store, "session")
    // Добавить AuthMiddleware, если нужно
    // r.Use(routing.AuthMiddleware(store, "session"))

    port := os.Getenv("PORT")
    if port == "" {
        port = "8083"
    }
    addr := ":" + port
    log.Printf("Starting server on %s", addr)
    if err := http.ListenAndServe(addr, csrfMw(r)); err != nil {
        log.Fatalf("Server error: %v", err)
    }
}