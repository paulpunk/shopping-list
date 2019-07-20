package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ant0ine/go-json-rest/rest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	mongoaddr, err := determineMongoDbAddress()
	if err != nil {
		log.Fatal(err)
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoaddr)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

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
