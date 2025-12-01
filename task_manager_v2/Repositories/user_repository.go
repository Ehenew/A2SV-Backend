package Repositories

import (
	"a2sv-backend/task_manager/Domain"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(user Domain.User) (Domain.User, error)
	FindByUsername(username string) (Domain.User, error)
	FindByID(id primitive.ObjectID) (Domain.User, error)
	Update(user Domain.User) (Domain.User, error)
	Count() (int64, error)
}

type mongoUserRepository struct {
	db *mongo.Database
}

func NewMongoUserRepository(db *mongo.Database) UserRepository {
	return &mongoUserRepository{db: db}
}

func (r *mongoUserRepository) Count() (int64, error) {
	return r.db.Collection("users").CountDocuments(context.Background(), bson.M{})
}

func (r *mongoUserRepository) Create(user Domain.User) (Domain.User, error) {
	_, err := r.db.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		return Domain.User{}, err
	}
	return user, nil
}

func (r *mongoUserRepository) FindByUsername(username string) (Domain.User, error) {
	var user Domain.User
	err := r.db.Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Domain.User{}, errors.New("user not found")
		}
		return Domain.User{}, err
	}
	return user, nil
}

func (r *mongoUserRepository) FindByID(id primitive.ObjectID) (Domain.User, error) {
	var user Domain.User
	err := r.db.Collection("users").FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Domain.User{}, errors.New("user not found")
		}
		return Domain.User{}, err
	}
	return user, nil
}

func (r *mongoUserRepository) Update(user Domain.User) (Domain.User, error) {
	_, err := r.db.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	if err != nil {
		return Domain.User{}, err
	}
	return user, nil
}
