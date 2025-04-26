package routing

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"otkritki/core/db"
	"otkritki/core/models"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"os"
)

func init() {
	dsn := fmt.Sprintf("root:%v@tcp(db:3306)/%v?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_PASSWORD"), os.Getenv("DBNAME"))
	database = db.NewDb(dsn, &models.User{}, &models.GiftCard{})
	cookieStore = sessions.NewCookieStore([]byte("sixtenbytelength"))
	encoder = schema.NewEncoder()
	decoder = schema.NewDecoder()
}

func abort(w http.ResponseWriter, errors ...string) {
	log.Println("Aborting...", errors)
	if len(errors) == 0 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		http.Error(w, strings.Join(errors, " "), http.StatusNotFound)
	}
}
