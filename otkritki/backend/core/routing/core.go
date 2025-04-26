package routing

import (
    "net/http"
    "os"
    "strings"
)

// CORSMiddleware – динамический CORS по списку из env
func CORSMiddleware(next http.Handler) http.Handler {
    allowed := strings.Split(os.Getenv("CORS_ALLOWED"), ",") // ИЗМЕНЕНО
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        for _, o := range allowed {
            if o == origin {
                w.Header().Set("Access-Control-Allow-Origin", origin) // ИЗМЕНЕНО
                break
            }
        }
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}
