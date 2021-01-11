package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
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

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)

	router := mux.NewRouter()

	router.HandleFunc("/people", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPeopleByID).Methods("GET")
	fmt.Println("Server on port:", 4000)
	http.ListenAndServe(":4000", router)

}

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

func GetPeople(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")

	var people []Person

	collection := client.Database("users_go").Collection("people")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message":"error"}`))
		return
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)

		people = append(people, person)
	}

	if err := cursor.Err(); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message":"error"}`))
		return
	}

	json.NewEncoder(res).Encode(people)
}

func GetPeopleByID(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	params := mux.Vars(req)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	var person Person

	collection := client.Database("users_go").Collection("people")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "error" }`))
		return
	}

	json.NewEncoder(res).Encode(person)
}
