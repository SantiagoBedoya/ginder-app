package mongodb

import (
	"context"
	"time"

	"github.com/SantiagoBedoya/ginder_users-api/users"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, err
}

func NewMongoRepository(mongoUrl, mongoDB string, mongoTimeout int) (users.UserRepository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}
	client, err := newMongoClient(mongoUrl, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepository")
	}
	repo.client = client
	return repo, nil
}

func (r *mongoRepository) FindAll() ([]users.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	list := make([]users.User, 0)
	collection := r.client.Database(r.database).Collection("users")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, "MongoRepository.User.FindAll")
	}
	for cursor.Next(ctx) {
		user := users.User{}
		if err := cursor.Decode(&user); err != nil {
			return nil, errors.Wrap(err, "MongoRepository.User.FindAll")
		}
		list = append(list, user)
	}
	return list, nil
}
func (r *mongoRepository) FindOneByID(id string) (*users.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	collection := r.client.Database(r.database).Collection("users")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "MongoRepository.User.FindOneById")
	}
	filter := bson.M{
		"_id": objectId,
	}
	user := &users.User{}
	if err := collection.FindOne(ctx, filter).Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(users.ErrUserNotFound, "MongoRepository.User.FindOneById")
		}
		return nil, errors.Wrap(err, "MongoRepository.User.FindOneById")
	}
	return user, nil
}

func (r *mongoRepository) FindOneByEmail(email string) (*users.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	collection := r.client.Database(r.database).Collection("users")
	filter := bson.M{
		"email": email,
	}
	user := &users.User{}
	if err := collection.FindOne(ctx, filter).Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(users.ErrUserNotFound, "MongoRepository.User.FindOneById")
		}
		return nil, errors.Wrap(err, "MongoRepository.User.FindOneById")
	}
	return user, nil
}

func (r *mongoRepository) Create(user *users.User) (*users.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	collection := r.client.Database(r.database).Collection("users")
	result, err := collection.InsertOne(ctx, bson.M{
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
		"password":  user.Password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "MongoRepository.User.Create")
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}
func (r *mongoRepository) UpdateOneByID(id string, user *users.User) error {
	_, err := r.FindOneByID(id)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	collection := r.client.Database(r.database).Collection("users")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, "MongoRepository.User.UpdateOneById")
	}
	filter := bson.M{
		"_id": objectId,
	}
	_, err = collection.UpdateOne(ctx, filter, bson.M{
		"$set": bson.M{
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
		},
	})
	if err != nil {
		return errors.Wrap(err, "MongoRepository.User.UpdateOneById")
	}
	return nil
}
func (r *mongoRepository) DeleteOneByID(id string) error {
	_, err := r.FindOneByID(id)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	collection := r.client.Database(r.database).Collection("users")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, "MongoRepository.User.DeleteOneByID")
	}
	filter := bson.M{
		"_id": objectId,
	}
	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "MongoRepository.User.DeleteOneByID")
	}
	return nil
}
