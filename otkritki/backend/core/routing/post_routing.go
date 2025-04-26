package routing

import (
    "encoding/json"
    "net/http"

    "github.com/IvanSt1/ctfad/otkritki/backend/core/db"
    "github.com/IvanSt1/ctfad/otkritki/backend/core/models"
)

// authUser сохраняет данные пользователя в сессии и возвращает их.
func authUser(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore, cookieName string, user *models.User) {
    session, _ := store.Get(r, cookieName)
    session.Values["authenticated"] = true
    session.Values["gender"] = string(user.Gender)
    session.Values["id"] = user.ID
    session.Save(r, w)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(user)
}

// RegisterPost навешивает POST-роуты на переданный mux.Router.
func RegisterPost(r *mux.Router, store *sessions.CookieStore, cookieName string) {
    // Регистрация
    r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Username string        `json:"username" schema:"username,required"`
            Password string        `json:"password" schema:"password,required"`
            Gender   models.Gender `json:"gender" schema:"gender,required"`
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
    }).Methods("POST")

    // Логин
    r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Username string `json:"username" schema:"username,required"`
            Password string `json:"password" schema:"password,required"`
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
    }).Methods("POST")

    // Добавление открытки
    r.HandleFunc("/api/cards", func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            To        string `json:"to" schema:"to,required"`
            Text      string `json:"text" schema:"text,required"`
            ImageType string `json:"imageType" schema:"imageType,required"`
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
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(card)
    }).Methods("POST")

    // Логаут
    r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, cookieName)
        session.Options.MaxAge = -1
        session.Save(r, w)
        w.WriteHeader(http.StatusNoContent)
    }).Methods("POST")
}
