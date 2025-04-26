package routing

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"

    "github.com/IvanSt1/ctfad/otkritki/backend/core/db"
)

// listCardsHandler возвращает JSON-список всех открыток.
func listCardsHandler(w http.ResponseWriter, r *http.Request) {
    cards, err := db.ListCards()
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(cards)
}

// RegisterRoutes навешивает GET-маршруты на переданный mux.Router.
func RegisterRoutes(r *mux.Router) {
    r.HandleFunc("/api/cards", listCardsHandler).Methods("GET")
}