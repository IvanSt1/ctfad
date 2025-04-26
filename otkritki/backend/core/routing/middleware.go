package routing

import (
    "net/http"

    // Абсолютный путь до вашего модуля и пакета моделей
    "github.com/IvanSt1/ctfad/otkritki/backend/core/models"
)

// CorsMiddleware выставляет CORS-заголовки
func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:31338")
        w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,Cookie,Accept,Pragma,Cache-Control,Expires")
        w.Header().Add("Access-Control-Allow-Credentials", "true")
        w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Add("Content-Type", "application/json;charset=UTF-8")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// AuthMiddleware проверяет сессию и существование пользователя
func AuthMiddleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        session, err := cookieStore.Get(r, cookieName)
        if err != nil || session.IsNew {
            abort(w, "Provide session cookie")
            return
        }
        auth, ok := session.Values[authKey].(bool)
        if !ok || !auth {
            abort(w, "Invalid session cookie")
            return
        }
        id, ok := session.Values[idKey].(uint)
        if !ok {
            abort(w, "Invalid session ID")
            return
        }
        if _, err := database.GetUserById(id); err != nil {
            abort(w, "Invalid user")
            return
        }
        h.ServeHTTP(w, r)
    })
}

// MaleWiddleWare пропускает только пользователей с gender == models.Male
func MaleWiddleWare(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        session, _ := cookieStore.Get(r, cookieName)
        gender, ok := session.Values[genderKey].(string)
        if !ok || gender != string(models.Male) {
            abort(w, "Only males are allowed")
            return
        }
        h.ServeHTTP(w, r)
    })
}
