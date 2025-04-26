package routing

import (
    "encoding/json"
    "log"
    "net/http"

    // Абсолютный путь вашего модуля
    "github.com/IvanSt1/ctfad/otkritki/backend/core/models"
)

// CheckAuth просто проверяет доступ (пригодится для тестов)
func CheckAuth(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}

// authUser сохраняет информацию в сессию и возвращает данные пользователя
func authUser(w http.ResponseWriter, r *http.Request, user *models.User) {
    session, _ := cookieStore.Get(r, cookieName)
    session.Values[authKey] = true
    session.Values[genderKey] = string(user.Gender)
    session.Values[idKey] = user.ID
    if err := session.Save(r, w); err != nil {
        abort(w, "Could not authenticate user")
        return
    }

    responseData, err := json.Marshal(user)
    if err != nil {
        abort(w, "Could not provide user data in return")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(responseData)
}

// RegisterPost обрабатывает регистрацию нового пользователя
func RegisterPost(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Username string        `json:"username" schema:"username,required"`
        Password string        `json:"password" schema:"password,required"`
        Gender   models.Gender `json:"gender" schema:"gender,required"`
    }
    if err := r.ParseForm(); err != nil {
        abort(w, "Could not parse Register params")
        return
    }
    if err := decoder.Decode(&req, r.PostForm); err != nil {
        abort(w, "Invalid Register params")
        return
    }
    if len(req.Username) < usernameLen || len(req.Password) < passwordLen {
        abort(w, "Username or password too short")
        return
    }

    user := &models.User{
        Username: req.Username,
        Password: req.Password,
        Gender:   req.Gender,
    }
    if _, err := database.AddUser(user); err != nil {
        abort(w, err.Error())
        return
    }
    authUser(w, r, user)
}

// LoginPost обрабатывает вход
func LoginPost(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Username string `json:"username" schema:"username,required"`
        Password string `json:"password" schema:"password,required"`
    }
    if err := r.ParseForm(); err != nil {
        abort(w, "Could not parse Login params")
        return
    }
    if err := decoder.Decode(&req, r.PostForm); err != nil {
        abort(w, "Invalid Login params")
        return
    }
    user, err := database.GetUserByName(req.Username)
    if err != nil {
        abort(w, err.Error())
        return
    }
    authUser(w, r, user)
}

// AddCardPost добавляет новую открытку (только для мужчин)
func AddCardPost(w http.ResponseWriter, r *http.Request) {
    var req struct {
        To        string `json:"to" schema:"to,required"`
        Text      string `json:"text" schema:"text,required"`
        ImageType string `json:"imageType" schema:"imageType,required"`
    }
    type CardResponse struct {
        *models.GiftCard
        Id uint `json:"id"`
    }

    if err := r.ParseForm(); err != nil {
        abort(w, "Invalid card parameters")
        return
    }
    if err := decoder.Decode(&req, r.PostForm); err != nil {
        abort(w, "Invalid card parameters")
        return
    }

    session, _ := cookieStore.Get(r, cookieName)
    userID := session.Values[idKey].(uint)
    sender, err := database.GetUserById(userID)
    if err != nil {
        abort(w, err.Error())
        return
    }

    newCard := &models.GiftCard{
        To:        req.To,
        From:      sender.Username,
        Text:      req.Text,
        ImageType: req.ImageType,
    }
    if _, err := database.AddCard(newCard); err != nil {
        abort(w, err.Error())
        return
    }

    response := &CardResponse{
        GiftCard: newCard,
        Id:       newCard.ID,
    }
    respData, err := json.Marshal(response)
    if err != nil {
        abort(w, "Could not serialize response")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(respData)
}

// LogoutPost обрабатывает выход из системы
func LogoutPost(w http.ResponseWriter, r *http.Request) {
    session, _ := cookieStore.Get(r, cookieName)
    session.Values[authKey] = false
    session.Options.MaxAge = -1
    if err := session.Save(r, w); err != nil {
        log.Println("ERROR logging out:", err)
    }
    w.WriteHeader(http.StatusNoContent)
}
