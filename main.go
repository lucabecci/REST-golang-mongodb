package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

var client *mongo.Client

func CreatePersonEndpoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	var person Person
	json.NewDecoder(req.Body).Decode(&person)

	collection := client.Database("users_go").Collection("people")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := collection.InsertOne(ctx, person)

	if err != nil {
		var errorPetition bytes.Buffer
		errorPetition.WriteString(`Error`)
		json.NewEncoder(res).Encode(errorPetition.String())
		return
	}

	json.NewEncoder(res).Encode(result)

	res.WriteHeader(http.StatusAccepted)
}

func main() {
	fmt.Printf("Hello World")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)

	fmt.Println("DB is connected")

	router := mux.NewRouter()

	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	http.ListenAndServe(":4000", router)

}
