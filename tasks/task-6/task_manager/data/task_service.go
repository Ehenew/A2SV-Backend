package data

import (
	"context"
	"errors"
	"log"
	"os"

	"a2sv-backend/task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

// initializing the MongoDB connection
func InitMongoDB() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI is not set in environment")
	}

	clientOptions := options.Client().ApplyURI(uri)
	
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB Atlas!")
	collection = client.Database("task_manager").Collection("tasks")
}

// GetAllTasks
func GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTaskByID
func GetTaskByID(id string) (*models.Task, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid task ID")
	}

	var task models.Task
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return &task, nil
}

// add a new task to the database
func AddTask(newTask models.Task) (primitive.ObjectID, error) {
	result, err := collection.InsertOne(context.TODO(), newTask)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// UpdateTask
func UpdateTask(id string, updatedTask models.Task) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID")
	}

	update := bson.M{}
	if updatedTask.Title != "" {
		update["title"] = updatedTask.Title
	}
	if updatedTask.Description != "" {
		update["description"] = updatedTask.Description
	}
	if !updatedTask.DueDate.IsZero() {
		update["duedate"] = updatedTask.DueDate
	}
	if updatedTask.Status != "" {
		update["status"] = updatedTask.Status
	}

	if len(update) == 0 {
		return nil
	}

	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}

// DeleteTask
func DeleteTask(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID")
	}

	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}
