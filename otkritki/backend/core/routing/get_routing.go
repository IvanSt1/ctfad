package routing

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"

    // Правильный путь вашего модуля
    "github.com/IvanSt1/ctfad/otkritki/backend/core/models"
)

// listCardsHandler возвращает JSON-список открыток
func listCardsHandler(w http.ResponseWriter, r *http.Request) {
    cards, err := models.ListCards()
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(cards)
}

// RegisterRoutes навешивает маршруты на переданный роутер
func RegisterRoutes(r *mux.Router) {
    r.HandleFunc("/api/cards", listCardsHandler).Methods("GET")
    // при необходимости добавьте POST, PUT, DELETE…
}
