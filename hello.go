package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionUsers *mongo.Collection
var collectionPosts *mongo.Collection
var ctx context.Context
var lock sync.Mutex

func homePage(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	var val string
	para := r.URL.Path
	if len(para) > 6 {
		val = para[7:]
	}
	if r.Method == "POST" && val == "" {
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
	} else if r.Method == "GET" && len(val) >= 1 {
		var data bson.M
		err := collectionUsers.FindOne(ctx, bson.M{"id": val}).Decode(&data)
		if err != nil {
			panic(err)
		} else {
			fmt.Fprintln(w, data)
		}
	}
	lock.Unlock()
}
func posts(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	var val string
	para := r.URL.Path
	if len(para) > 6 {
		val = para[7:]
	}

	if r.Method == "POST" && val == "" {
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
	} else if r.Method == "GET" && len(val) >= 1 {

		var data bson.M
		err := collectionPosts.FindOne(ctx, bson.M{"id": val}).Decode(&data)
		if err != nil {
			panic(err)
		} else {
			fmt.Fprintln(w, data)
		}
	}
	lock.Unlock()
}
func getAllPosts(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	option := options.Find()
	option.SetLimit(2)
	var val string
	para := r.URL.Path
	if len(para) > 6 {
		val = para[13:]
		fmt.Println(val)
	}
	if r.Method == "GET" && val != "" {
		var data []bson.M
		cursor, err := collectionPosts.Find(ctx, bson.M{"userId": val}, option)
		if err != nil {
			panic(err)
		} else {
			err = cursor.All(ctx, &data)
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(w, data)
		}
	}
	lock.Unlock()
}
func TestPostingUserDetails(t *testing.T){
	Body := strings.NewReader({"id": "1234", "name": "Aditya Keshan", 
	"email": "aditya.keshan2019@vitstudent.ac.in", "password": "ABcd123@@"})
		req, err := http.NewRequest("POST", "http://localhost:3000/users/", Body)
		if(err!=nil){
			fmt.Println("Error in Posting");
		}else{
			fmt.Println("Success in testing:",req.Close);
		}
}
func TestPostingPostsDetails(t *testing.T){
	
	Body := strings.NewReader({"id": "1234", "caption": "Be the change", 
		"url": "https://www.google.com", "timestamp": "16-47","userId":"1234"})
	req, err := http.NewRequest("POST", "http://localhost:3000/posts/", Body)
	if(err!=nil){
		fmt.Println("Error in Posting");
	}else{
		fmt.Println("Success in testing:",req.Close);
	}
}
func TestGetUserDetails(t *testing.T)
{
	req, err := http.NewRequest("GET", "http://localhost:3000/users/1234", nil)
	if(err!=nil){
		fmt.Println("Error in Get");
	}else{
		fmt.Println("Success in testing:",req.Close);
	}
}
func TestGetPostDetails(t *testing.T)
{
	req, err := http.NewRequest("GET", "http://localhost:3000/posts/1234", nil)
	if(err!=nil){
		fmt.Println("Error in Get");
	}else{
		fmt.Println("Success in testing:",req.Close);
	}
}
func TestGetUsersAllPostsDetails(t *testing.T)
{
	req, err := http.NewRequest("GET", "http://localhost:3000/posts/users/1234", nil)
	if(err!=nil){
		fmt.Println("Error in Get");
	}else{
		fmt.Println("Success in testing:",req.Close);
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
