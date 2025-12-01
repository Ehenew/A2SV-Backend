package Repositories

import (
	"a2sv-backend/task_manager/Domain"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository interface {
	Create(task Domain.Task) (Domain.Task, error)
	FindAll() ([]Domain.Task, error)
	FindByID(id primitive.ObjectID) (Domain.Task, error)
	Update(task Domain.Task) (Domain.Task, error)
	Delete(id primitive.ObjectID) error
}

type mongoTaskRepository struct {
	db *mongo.Database
}

func NewMongoTaskRepository(db *mongo.Database) TaskRepository {
	return &mongoTaskRepository{db: db}
}

func (r *mongoTaskRepository) Create(task Domain.Task) (Domain.Task, error) {
	result, err := r.db.Collection("tasks").InsertOne(context.Background(), task)
	if err != nil {
		return Domain.Task{}, err
	}
	task.ID = result.InsertedID.(primitive.ObjectID)
	return task, nil
}

func (r *mongoTaskRepository) FindAll() ([]Domain.Task, error) {
	cursor, err := r.db.Collection("tasks").Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var tasks []Domain.Task
	if err = cursor.All(context.Background(), &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *mongoTaskRepository) FindByID(id primitive.ObjectID) (Domain.Task, error) {
	var task Domain.Task
	err := r.db.Collection("tasks").FindOne(context.Background(), bson.M{"_id": id}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Domain.Task{}, errors.New("task not found")
		}
		return Domain.Task{}, err
	}
	return task, nil
}

func (r *mongoTaskRepository) Update(task Domain.Task) (Domain.Task, error) {
	_, err := r.db.Collection("tasks").UpdateOne(
		context.Background(),
		bson.M{"_id": task.ID},
		bson.M{"$set": task},
	)
	if err != nil {
		return Domain.Task{}, err
	}
	return task, nil
}

func (r *mongoTaskRepository) Delete(id primitive.ObjectID) error {
	_, err := r.db.Collection("tasks").DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
