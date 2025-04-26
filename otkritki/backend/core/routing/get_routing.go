package routing

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/IvanSt1/ctfad/otkritki/backend/core/db"
)

// RegisterRoutes настраивает GET-маршруты
func RegisterRoutes(r *mux.Router) {
    r.HandleFunc("/api/cards", listCardsHandler).Methods("GET")
}

// listCardsHandler возвращает все открытки в формате JSON
func listCardsHandler(w http.ResponseWriter, r *http.Request) {
    cards, err := db.ListCards()
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cards)
}