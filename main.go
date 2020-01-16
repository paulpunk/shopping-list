package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ant0ine/go-json-rest/rest"
)

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func determineMongoDbAddress() (string, error) {
	address := os.Getenv("MONGODB_URI")
	if address == "" {
		return "", fmt.Errorf("$MONGODB_URI not set")
	}
	return address, nil
}
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello Katy")
}

func main() {
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/list", sync),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	log.Printf("Listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, api.MakeHandler()); err != nil {
		panic(err)
	}

}

func sync(w rest.ResponseWriter, r *rest.Request) {

	nicelist := Nicelist{}
	err := r.DecodeJsonPayload(&nicelist)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if nicelist.User == "" {
		rest.Error(w, "user required", 400)
		return
	}

	persist(&nicelist)

	w.WriteJson(&nicelist)
}
