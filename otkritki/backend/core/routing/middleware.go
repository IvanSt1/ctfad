package routing

import (
	"otkritki/core/models"
	"net/http"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:31338")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,Cookie,Accept,Pragma,Cache-Control,Expires")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Add("Content-Type", "application/json;charset=UTF-8")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := cookieStore.Get(r, cookieName)
		if err != nil || session.IsNew {
			abort(w, "Provide session cookie")
			return
		}
		auth := session.Values[authKey].(bool)
		id := session.Values[idKey].(uint)

		if !auth {
			abort(w, "Invalid session cookie")
			return
		} else if _, err := database.GetUserById(id); err != nil {
			abort(w, "Invalid user")
			return
		}
		h.ServeHTTP(w, r)
	})
}

func MaleWiddleWare(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := cookieStore.Get(r, cookieName)
		if gender, ok := session.Values[genderKey].(string); !ok || gender != string(models.Male) {
			abort(w, "Only males are allowed")
			return
		}
		h.ServeHTTP(w, r)
	})
}
