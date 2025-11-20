package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global variable to hold the MongoDB client
var client *mongo.Client

// Trainer struct defines the data model for our collection
type Trainer struct {
	// Use primitive.ObjectID for MongoDB's unique ID field
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name,omitempty"`
	Age  int                `bson:"age,omitempty"`
	City string             `bson:"city,omitempty"`
}

// Global collection handle
var collection *mongo.Collection

func main() {
	// =======================================
	// 1. ATLAS CONNECTION SETUP
	// =======================================
	uri := "mongodb+srv://ehenewamogne:uMNQ7U5iLBLSdBAl@cluster0.y8khhhs.mongodb.net/?appName=Cluster0"

	// Create a context with a 10-second timeout for the initial connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	var err error

	// Connect to MongoDB Atlas
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB Atlas:", err)
	}

	// Ping to verify the connection
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Could not verify connection to Atlas:", err)
	}

	fmt.Println("Successfully connected to MongoDB Atlas!")

	// Get a handle to the 'test' database and 'trainers' collection
	collection = client.Database("test").Collection("trainers")

	// =======================================
	// 2. CRUD OPERATIONS EXECUTION
	// =======================================
	fmt.Println("\n--- Starting CRUD Operations ---")

	// Ensure a clean slate for the demo
	deleteManyDocuments(bson.D{})

	// 1. Insert Operations
	ashID := insertOneDocument()
	insertManyDocuments()

	// 2. Find Operations
	findSingleDocument(ashID)
	findAllDocuments()

	// 3. Update Operations
	updateOneDocument(ashID)
	findSingleDocument(ashID) // Show the update took effect

	// 4. Delete Operations (Cleanup)
	// deleteOneDocument(ashID)
	deleteManyDocuments(bson.D{})

	fmt.Println("\n--- CRUD Operations Complete ---")

	// Disconnect the client when the program exits
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal("Error disconnecting:", err)
		}
		fmt.Println("Connection to MongoDB Atlas closed.")
	}()
}

// --- CRUD Function Implementations ---

// Insert a single Trainer document
func insertOneDocument() primitive.ObjectID {
	ash := Trainer{Name: "Ash", Age: 10, City: "Pallet Town"}

	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("1. Inserted a single document. ID: %v\n", insertResult.InsertedID)
	// Return the new ID for later use
	return insertResult.InsertedID.(primitive.ObjectID)
}

// Insert multiple Trainer documents
func insertManyDocuments() {
	misty := Trainer{Name: "Misty", Age: 10, City: "Cerulean City"}
	brock := Trainer{Name: "Brock", Age: 15, City: "Pewter City"}

	// InsertMany requires a slice of interface{}
	trainers := []interface{}{misty, brock}

	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("1. Inserted multiple documents. Count: %d\n", len(insertManyResult.InsertedIDs))
}

// Find a single Trainer document by ID
func findSingleDocument(id primitive.ObjectID) {
	var result Trainer

	// Filter by the unique ObjectID
	filter := bson.D{{"_id", id}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("2. FindOne: No document found with ID %v\n", id)
		return
	} else if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("2. Found single document: Name=%s, Age=%d, City=%s\n", result.Name, result.Age, result.City)
}

// Find all documents and iterate through the cursor
func findAllDocuments() {
	// Passing an empty filter bson.D{} matches all documents
	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.TODO())

	var results []Trainer

	// Iterate through the cursor and decode each document
	for cur.Next(context.TODO()) {
		var elem Trainer
		if err := cur.Decode(&elem); err != nil {
			log.Fatal(err)
		}
		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("2. Found %d documents in total (Misty, Brock, and Ash):\n", len(results))
	for _, t := range results {
		fmt.Printf("   - %s\n", t.Name)
	}
}

// Update a single document
func updateOneDocument(id primitive.ObjectID) {
	// Filter to match the document by ID
	filter := bson.D{{"_id", id}}

	// Update document: use $set to change a field, $inc to increment
	update := bson.D{
		{"$set", bson.D{{"city", "New Pallet Town"}}},
		{"$inc", bson.D{{"age", 1}}}, // Happy Birthday Ash!
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("3. Updated %v document(s). New City/Age for Ash.\n", updateResult.ModifiedCount)
}

// Delete all documents matching the filter (in this case, all documents)
func deleteManyDocuments(filter interface{}) {
	// Passing bson.D{} as the filter matches all documents
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("4. Deleted %v documents (Cleanup).\n", deleteResult.DeletedCount)
}
