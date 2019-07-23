package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ant0ine/go-json-rest/rest"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func init() {
	mongoaddr, err := determineMongoDbAddress()
	if err != nil {
		log.Fatal(err)
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoaddr)

	// Connect to MongoDB
	Client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
}

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

type List struct {
	User  string
	Items []*Item
}

type Item struct {
	ID      int
	Version int
	Name    string
	Checked bool
	User    string
	List    string
	State   string
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

	list := List{}
	err := r.DecodeJsonPayload(&list)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if list.User == "" {
		rest.Error(w, "user required", 400)
		return
	}

	fmt.Println("Connected to MongoDB!")

	collection := Client.Database("heroku_tx1qdrzx").Collection("item")

	for _, item := range list.Items {
		if item.State == "create" {
			create(collection, item)
		}
		if item.State == "update" {
			update(collection, item)
		}
		if item.State == "delete" {
			delete(collection, item)
		}
	}

	list.Items = find(collection, &list.User)

	w.WriteJson(&list)
}

func create(collection *mongo.Collection, item *Item) {

	insertResult, err := collection.InsertOne(context.TODO(), bson.M{
		"id":      item.ID,
		"version": item.Version,
		"user":    item.User,
		"name":    item.Name,
		"checked": item.Checked,
		"list":    item.List,
	},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func update(collection *mongo.Collection, item *Item) {
	filter := bson.D{{"name", "Ash"}}

	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func delete(collection *mongo.Collection, item *Item) {

	filter := bson.D{{"id", item.ID}, {"user", item.User}}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}

func find(collection *mongo.Collection, user *string) []*Item {
	// Here's an array in which you can store the decoded documents
	var results []*Item

	// Passing bson.D{{}} as the filter matches all documents in the collection
	filter := bson.D{{"user", user}}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Item
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	return results
}
