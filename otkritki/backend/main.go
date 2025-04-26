package main

import (
    "encoding/json"
    "net/http"
    "os"
    "strings"

    "github.com/gorilla/csrf"
    "github.com/gorilla/mux"
    "github.com/gorilla/sessions"
    "golang.org/x/time/rate"

    "github.com/IvanSt1/ctfad/backend/core/models"
    "github.com/IvanSt1/ctfad/backend/core/routing"
)

var (
    store       *sessions.CookieStore
    loginLimiter *rate.Limiter                // ИЗМЕНЕНО: для rate-limiting
)

func frontPorch(w http.ResponseWriter, r *http.Request) {
    // этот «бекдор» полностью убран
    http.Error(w, "Forbidden", http.StatusForbidden)    // ИЗМЕНЕНО
}

func main() {
    // === Настройка CookieStore ===
    authKey := []byte(os.Getenv("COOKIE_AUTH_KEY"))     // ИЗМЕНЕНО
    encKey  := []byte(os.Getenv("COOKIE_ENC_KEY"))      // ИЗМЕНЕНО
    store = sessions.NewCookieStore(authKey, encKey)    // ИЗМЕНЕНО
    store.Options = &sessions.Options{                   // ИЗМЕНЕНО
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
    }

    // === Rate limiter для /login и /register ===
    loginLimiter = rate.NewLimiter(1, 5)                // 1 запрос в секунду, буфер 5

    // === Маршруты и middlewares ===
    r := mux.NewRouter()
    // Динамический CORS
    r.Use(routing.CORSMiddleware)                       // ИЗМЕНЕНО
    // CSRF
    csrfMiddleware := csrf.Protect(
        []byte(os.Getenv("CSRF_KEY")),                  // ИЗМЕНЕНО
        csrf.Secure(true),
        csrf.Path("/"),
    )
    // Auth handlers
    r.HandleFunc("/login", withRateLimit(loginHandler)).Methods("POST")       // ИЗМЕНЕНО
    r.HandleFunc("/register", withRateLimit(registerHandler)).Methods("POST") // ИЗМЕНЕНО

    // Прочие API
    r.HandleFunc("/api/cards", createCardHandler).Methods("POST")
    r.HandleFunc("/api/cards", listCardsHandler).Methods("GET")

    // Убираем «бекдор» на фронтпорч
    //r.HandleFunc("/api/nothingtoseehere", frontPorch).Methods("GET")

    http.ListenAndServe(":8080", csrfMiddleware(r))
}

// withRateLimit оборачивает handler в rate-limiter
func withRateLimit(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !loginLimiter.Allow() {
            http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
            return
        }
        next(w, r)
    }
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    type creds struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    var c creds
    if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }
    user, err := models.FindUserByUsername(c.Username)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }
    // bcrypt.compare
    if err := models.ComparePassword(user.PasswordHash, c.Password); err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }
    // Успешный логин
    session, _ := store.Get(r, "session-name")
    session.Values["user_id"] = user.ID
    session.Save(r, w)

    w.Header().Set("Content-Type", "application/json") // ИЗМЕНЕНО: заголовок
    w.WriteHeader(http.StatusOK)                       // ИЗМЕНЕНО: WriteHeader до Write
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    type req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    var data req
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }
    // валидация и bcrypt-хеш
    if err := models.CreateUser(data.Username, data.Password); err != nil {
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json") // ИЗМЕНЕНО
    w.WriteHeader(http.StatusCreated)                   // ИЗМЕНЕНО
    json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}
