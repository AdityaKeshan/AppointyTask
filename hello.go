package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionUsers *mongo.Collection
var collectionPosts *mongo.Collection
var ctx context.Context

func homePage(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if !ok {
		if r.Method == "POST" && ids == nil {
			id := r.FormValue("id")
			name := r.FormValue("name")
			email := r.FormValue("email")
			password := r.FormValue("password")
			var document interface{}
			document = bson.D{
				{"id", id}, {"name", name}, {"email", email}, {"password", password}}

			res, err := collectionUsers.InsertOne(ctx, document)
			if err != nil {
				panic(err)
			}
			fmt.Println(res.InsertedID)
		}
	} else {
		if r.Method == "GET" && len(ids[0]) >= 1 {

			key := ids[0]
			var podcast bson.M
			err := collectionUsers.FindOne(ctx, bson.M{"id": key}).Decode(&podcast)
			if err != nil {
				panic(err)
			} else {
				fmt.Println(podcast["name"])
			}

		}
	}
}
func posts(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if !ok {
		if r.Method == "POST" && ids == nil {
			id := r.FormValue("id")
			caption := r.FormValue("caption")
			url := r.FormValue("url")
			timestamp := r.FormValue("timestamp")
			var document interface{}
			document = bson.D{
				{"id", id}, {"caption", caption}, {"url", url}, {"timestamp", timestamp}}

			res, err := collectionPosts.InsertOne(ctx, document)
			if err != nil {
				panic(err)
			}
			fmt.Println(res.InsertedID)
		}
	} else {
		if r.Method == "GET" && ids != nil {
			key := ids[0]
			var podcast bson.M
			err := collectionPosts.FindOne(ctx, bson.M{"id": key}).Decode(&podcast)
			if err != nil {
				panic(err)
			} else {
				fmt.Println(podcast["caption"])
			}
		}
	}
}
func getAllPosts(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if ok && r.Method == "GET" && ids != nil {

	}
}
func handleRequests() {
	http.HandleFunc("/users/", homePage)
	http.HandleFunc("/posts/", posts)
	http.HandleFunc("/posts/users/", getAllPosts)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err == nil {
		fmt.Println("Client connect starting....")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err == nil {
		fmt.Println("Connected")
	}
	collectionUsers = client.Database("userDB").Collection("user")
	collectionPosts = client.Database("userDB").Collection("post")
	handleRequests()

}
