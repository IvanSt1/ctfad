package routing

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/gorilla/sessions"

    "github.com/IvanSt1/ctfad/otkritki/backend/core/db"
    "github.com/IvanSt1/ctfad/otkritki/backend/core/models"
)

// RegisterPost настраивает POST-маршруты
func RegisterPost(r *mux.Router, store *sessions.CookieStore, cookieName string) {
    r.HandleFunc("/register", registerHandler(store, cookieName)).Methods("POST")
    r.HandleFunc("/login", loginHandler(store, cookieName)).Methods("POST")
    r.HandleFunc("/api/cards", addCardHandler(store, cookieName)).Methods("POST")
    r.HandleFunc("/logout", logoutHandler(store, cookieName)).Methods("POST")
}

func registerHandler(store *sessions.CookieStore, cookieName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Username string        `json:"username"`
            Password string        `json:"password"`
            Gender   models.Gender `json:"gender"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            abort(w, "Invalid request body")
            return
        }
        user := &models.User{Username: req.Username, PasswordHash: req.Password, Gender: req.Gender}
        created, err := db.AddUser(user)
        if err != nil {
            abort(w, err.Error())
            return
        }
        authUser(w, r, store, cookieName, created)
    }
}

func loginHandler(store *sessions.CookieStore, cookieName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            abort(w, "Invalid request body")
            return
        }
        user, err := db.GetUserByName(req.Username)
        if err != nil {
            abort(w, err.Error())
            return
        }
        authUser(w, r, store, cookieName, user)
    }
}

func addCardHandler(store *sessions.CookieStore, cookieName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            To        string `json:"to"`
            Text      string `json:"text"`
            ImageType string `json:"imageType"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            abort(w, "Invalid request body")
            return
        }
        session, _ := store.Get(r, cookieName)
        uid := session.Values["id"].(uint)
        sender, err := db.GetUserById(uid)
        if err != nil {
            abort(w, err.Error())
            return
        }
        card := &models.GiftCard{To: req.To, From: sender.Username, Text: req.Text, ImageType: req.ImageType}
        if err := db.CreateCard(card); err != nil {
            abort(w, err.Error())
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(card)
    }
}

func logoutHandler(store *sessions.CookieStore, cookieName string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, cookieName)
        session.Options.MaxAge = -1
        session.Save(r, w)
        w.WriteHeader(http.StatusNoContent)
    }
}