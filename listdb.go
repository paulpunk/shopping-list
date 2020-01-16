package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func createList(collection *mongo.Collection, list *List) {

	insertResult, err := collection.InsertOne(context.TODO(), bson.M{
		"user":       list.User,
		"name":       list.Name,
		"sharedwith": list.SharedWith,
	},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func updateList(collection *mongo.Collection, list *List) {
	//TODO: merge data in case of version conflict?

	filter := bson.D{{"id", item.ID}, {"user", item.User}}

	update := bson.D{
		{"$inc", bson.D{
			{"version", 1},
		}},
		{"$set", bson.D{
			{"name", item.Name},
		}},
		{"$set", bson.D{
			{"checked", item.Checked},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func deleteList(collection *mongo.Collection, list *List) {

	filter := bson.D{{"id", item.ID}, {"user", item.User}}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}

func findLists(user *string) []*List {
	collection := Client.Database("heroku_tx1qdrzx").Collection("item")

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
