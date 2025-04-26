package routing

import (
	"encoding/json"
	"otkritki/core/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func CardPage(w http.ResponseWriter, r *http.Request) {
	cardId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		abort(w, "Invalid card Id specified")
		return
	}
	card, err := database.GetCardById(uint(cardId))
	if err != nil {
		abort(w, "No such card found")
		return
	}
	session, _ := cookieStore.Get(r, cookieName)
	userID := session.Values[idKey].(uint)
	user, _ := database.GetUserById(userID)

	if !strings.Contains(card.To, user.Username) && card.From != user.Username {
		abort(w, "This card does not belong to you")
		return
	}

	response, _ := json.Marshal(&card)

	w.WriteHeader(200)
	if _, err := w.Write(response); err != nil {
		log.Println(err)
	}
}

func GetCards(w http.ResponseWriter, r *http.Request) {
	var response []models.GiftCard

	session, _ := cookieStore.Get(r, cookieName)
	userID := session.Values[idKey].(uint)
	user, _ := database.GetUserById(userID)

	if latestVar := r.URL.Query().Get("latest"); latestVar == "" {
		recvCards, _ := database.GetCardByReceiver(user.Username)
		response = append(response, *recvCards...)
	} else if latest, err := strconv.Atoi(latestVar); err != nil {
		abort(w, "Invalid latest variable")
        return
	} else {
		latestCards, _ := database.GetLatestCards(user.Username, latest)
		response = append(response, *latestCards...)
	}

	resp, _ := json.Marshal(&response)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		log.Println(err)
	}
}
