package routing

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/sessions"
    "github.com/IvanSt1/ctfad/otkritki/backend/core/db"
)

// abort отвечает JSON-ошибкой и статусом Bad Request
func abort(w http.ResponseWriter, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// CorsMiddleware устанавливает CORS-заголовки
func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// AuthMiddleware проверяет сессию и существование пользователя
func AuthMiddleware(store *sessions.CookieStore, cookieName string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            session, err := store.Get(r, cookieName)
            if err != nil || session.IsNew {
                abort(w, "Provide session cookie")
                return
            }
            auth, ok := session.Values["authenticated"].(bool)
            if !ok || !auth {
                abort(w, "Invalid session cookie")
                return
            }
            id, ok := session.Values["id"].(uint)
            if !ok {
                abort(w, "Invalid session ID")
                return
            }
            if _, err := db.GetUserById(id); err != nil {
                abort(w, "Invalid user")
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}