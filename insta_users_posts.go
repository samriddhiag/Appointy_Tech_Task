package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() { //Database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ := mongo.NewClient(options.Client().ApplyURI(secretDbURI))

	err := client.Connect(ctx)

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err != nil {
		log.Fatal("Couldn't connect to database", err)
	} else {
		log.Println("Connected to database !")
	}

	r := NewRouter()
	r.Methods(http.MethodGet).Handler(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Dummy response from Instagram API\n")
	}))

	
	//Create an user - POST request
	r.Methods(http.MethodPost).Handler(`/users`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		var user User
		json.NewDecoder(r.Body).Decode(&user)
		collection := client.Database("instadb").Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var userName string = user.Name
		encryptedPassword := hex.EncodeToString(encrypt([]byte(user.Password), secretHashKey))
		var userEmail string = user.Email

		println(userName, userEmail, encryptedPassword)

		doc := bson.M{"name": userName, "email": userEmail, "password": encryptedPassword}
		result, err := collection.InsertOne(ctx, doc)
		if err != nil {
			fmt.Fprint(w, "Error Creating Post !\n", err)
		} else {
			fmt.Fprint(w, "User Created !\n", result)
		}
	}))
	
	//Create a post - POST request
	r.Methods(http.MethodPost).Handler(`/posts`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		var post Posts
		json.NewDecoder(r.Body).Decode(&post)
		collection := client.Database("instadb").Collection("posts")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		doc := bson.M{"caption": post.Caption, "url": post.Url, "currentTime": post.CurrentTime, "userID": post.UserID}
		result, err := collection.InsertOne(ctx, doc)
		if err != nil {
			fmt.Fprint(w, "Error Creating Post !\n", err)
		} else {
			fmt.Fprint(w, "Post Created !\n", result)
		}
	}))

	
	//Get a user using id - GET request
	r.Methods(http.MethodGet).Handler(`/posts/users/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Pagination - restricted to 10 posts at a time
		r.ParseForm()
		var limit int64 = 10
		var skip int64 = 0

		if len(r.Form) > 0 {
			for k, v := range r.Form {
				if k == "skip" {
					skip, _ = strconv.ParseInt(v[0], 10, 64)
				}
				if k == "limit" {
					limit, _ = strconv.ParseInt(v[0], 10, 64)
				}
			}
		}

		multiOptions := options.Find().SetLimit(limit).SetSkip(skip)

		id := GetParam(r.Context(), "id")
		collection := client.Database("instadb").Collection("posts")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.M{"userID": id}

		var result []bson.M

		findCursor, err := collection.Find(ctx, filter, multiOptions)
		findCursor.All(ctx, &result)
		j, _ := json.Marshal(result)

		if err != nil {
			fmt.Fprint(w, "Error finding Posts !\n", err)
		} else {
			fmt.Fprint(w, "Found All Posts !\n", string(j))
		}

	}))

	//List all posts of a user - GET request
	r.Methods(http.MethodGet).Handler(`/users/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		collection := client.Database("instadb").Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		oid, _ := primitive.ObjectIDFromHex(id)

		filter := bson.M{"_id": oid}

		var result bson.M

		err := collection.FindOne(ctx, filter).Decode(&result)
		j, _ := json.Marshal(result)

		// var decryptedPassword string

		// fmt.Println(decryptedPassword)
		var user User
		json.Unmarshal([]byte(j), &user)
		hexPass, _ := hex.DecodeString(user.Password)
		decryptedPassword := string(decrypt(hexPass, secretHashKey))
		user.Password = decryptedPassword
		final, _ := json.Marshal(user)

		if err != nil {
			fmt.Fprint(w, "Error finding User !\n", err)
		} else {
			fmt.Fprint(w, "User Found !\n", string(final))
		}

	}))

	//Get a post using id - GET request
	r.Methods(http.MethodGet).Handler(`/posts/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		collection := client.Database("instadb").Collection("posts")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		oid, _ := primitive.ObjectIDFromHex(id)

		filter := bson.M{"_id": oid}

		var result bson.M
		err := collection.FindOne(ctx, filter).Decode(&result)
		j, _ := json.Marshal(result)
		if err != nil {
			fmt.Fprint(w, "Error finding Post !\n", err)
		} else {
			fmt.Fprint(w, "Post Found !\n", string(j))
		}
	}))

	http.ListenAndServe(secretPort, r)
	fmt.Println("Server listening on port : ", secretPort)
}