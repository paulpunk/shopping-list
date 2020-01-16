package main

import (
	"context"
	"fmt"
	"log"

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

	fmt.Println("Connected to MongoDB!")
}

func persist(nicelist *Nicelist) {

	listcollection := Client.Database("heroku_tx1qdrzx").Collection("list")
	itemcollection := Client.Database("heroku_tx1qdrzx").Collection("item")

	for _, list := range nicelist.Lists {
		if list.State == "create" {
			createList(listcollection, list)
		}
		if list.State == "update" {
			updateList(listcollection, list)
		}
		if list.State == "delete" {
			deleteList(listcollection, list)
		}
	}

	nicelist.Lists = findLists(&nicelist.User)

	for _, item := range nicelist.Items {
		if item.State == "create" {
			createItem(itemcollection, item)
		}
		if item.State == "update" {
			updateItem(itemcollection, item)
		}
		if item.State == "delete" {
			deleteItem(itemcollection, item)
		}
	}

	nicelist.Items = findItems(&nicelist.User)
}
