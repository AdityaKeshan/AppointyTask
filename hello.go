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
			var data bson.M
			err := collectionUsers.FindOne(ctx, bson.M{"id": key}).Decode(&data)
			if err != nil {
				panic(err)
			} else {
				fmt.Println(data["name"])
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
			userId := r.FormValue("userId")
			var document interface{}
			document = bson.D{
				{"id", id}, {"caption", caption}, {"url", url}, {"timestamp", timestamp}, {"userId", userId}}

			res, err := collectionPosts.InsertOne(ctx, document)
			if err != nil {
				panic(err)
			}
			fmt.Println(res.InsertedID)
		}
	} else {
		if r.Method == "GET" && ids != nil {
			key := ids[0]
			var data bson.M
			err := collectionPosts.FindOne(ctx, bson.M{"id": key}).Decode(&data)
			if err != nil {
				panic(err)
			} else {
				fmt.Println(data["caption"])
			}
		}
	}
}
func getAllPosts(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if ok && r.Method == "GET" && ids != nil {
		key := ids[0]
		var data []bson.M
		cursor, err := collectionPosts.Find(ctx, bson.M{"userId": key})
		if err != nil {
			panic(err)
		} else {
			err = cursor.All(ctx, &data)
			if err != nil {
				panic(err)
			}
			fmt.Println(data)
		}
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
