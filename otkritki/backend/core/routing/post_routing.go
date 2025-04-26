package routing

import (
	"encoding/json"
    "log"
	"otkritki/core/models"
	"net/http"
)

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func authUser(w http.ResponseWriter, r *http.Request, user *models.User) {
	session, _ := cookieStore.Get(r, cookieName)
	session.Values[authKey] = true
	session.Values[genderKey] = string(user.Gender)
	session.Values[idKey] = user.ID

	if err := session.Save(r, w); err != nil {
		abort(w, "Could not authenticate user")
        return
	}
	if responseData, err := json.Marshal(user); err != nil {
		abort(w, "Could not provide user data in return")
        return
	} else {
		w.Write(responseData)
		w.WriteHeader(http.StatusOK)
	}
}

func RegisterPost(w http.ResponseWriter, request *http.Request) {
	var RegisterRequest struct {
		Username string      `json:"username" schema:"username,required"`
		Password string      `json:"password" schema:"password,required"`
		Gender   models.Gender `json:"gender" schema:"gender,required"`
	}

	if err := request.ParseForm(); err != nil {
		abort(w, "Could not parse Register params")
		return
	}
	if err := decoder.Decode(&RegisterRequest, request.PostForm); err != nil {
		abort(w, "Invalid Register params")
		return
	}
	if len(RegisterRequest.Username) < usernameLen || len(RegisterRequest.Password) < passwordLen {
		abort(w, "Username or password too short")
		return
	}

	user := &models.User{
		Username: RegisterRequest.Username,
		Password: RegisterRequest.Password,
		Gender:   RegisterRequest.Gender,
	}

	if _, err := database.AddUser(user); err != nil {
		abort(w, err.Error())
		return
	} else {
		authUser(w, request, user)
	}
}

func LoginPost(w http.ResponseWriter, request *http.Request) {
	var LoginRequest struct {
		Username string `json:"username" schema:"username,required"`
		Password string `json:"password" validate:"password,required"`
	}

	if err := request.ParseForm(); err != nil {
		abort(w, "Could not parse Login params")
		return
	}
	if err := decoder.Decode(&LoginRequest, request.PostForm); err != nil {
		abort(w, "Invalid Login params")
		return
	}

	if user, err := database.GetUserByName(LoginRequest.Username); err == nil {
		authUser(w, request, user)
	} else {
		abort(w, err.Error())
	}

}

// MALE ONLY
func AddCardPost(w http.ResponseWriter, request *http.Request) {
	var CardRequest struct {
		To        string `json:"to" schema:"to,required"`
		Text      string `json:"text" schema:"text,required"`
		ImageType string `json:"imageType" schema:"imageType,required"`
	}

	type CardResponse struct {
		*models.GiftCard
		Id uint `json:"id"`
	}

	if err := request.ParseForm(); err != nil {
		abort(w, "Invalid card parameters")
		return
	}
	if err := decoder.Decode(&CardRequest, request.PostForm); err != nil {
		abort(w, "Invalid card parameters")
		return
	}

	session, _ := cookieStore.Get(request, cookieName)
	userID := session.Values[idKey].(uint)
	sender, err := database.GetUserById(userID)
	if err != nil {
		abort(w, err.Error())
		return
	}
	newCard := &models.GiftCard{
		To:        CardRequest.To,
		From:      sender.Username,
		Text:      CardRequest.Text,
		ImageType: CardRequest.ImageType,
	}
	if _, err := database.AddCard(newCard); err != nil {
		abort(w, err.Error())
		return
	}

	response := &CardResponse{
		newCard,
		newCard.ID,
	}

	if resp, err := json.Marshal(response); err == nil {
		w.Write(resp)
		w.WriteHeader(http.StatusOK)
	}
}

func LogoutPost(w http.ResponseWriter, request *http.Request) {
	session, _ := cookieStore.Get(request, cookieName)
	session.Values[authKey] = false
	session.Options.MaxAge = -1
	if err := session.Save(request, w); err != nil {
		log.Println("ERROR loggin out: ", err)
	}
}
