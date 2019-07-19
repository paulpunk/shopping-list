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
	api.SetApp(rest.AppSimple(func(w rest.ResponseWriter, r *rest.Request) {
		w.WriteJson(map[string]string{"Body": "Hello World!"})
	}))

	log.Printf("Listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, api.MakeHandler()); err != nil {
		panic(err)
	}



}
