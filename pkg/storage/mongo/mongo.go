package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Store - структура для работы с MongoDB.
type Store struct {
	db *mongo.Collection
}

// New создает новый экземпляр Store.
func New(connectionString string) (*Store, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database("gonews").Collection("posts") // Используем базу данных "gonews" и коллекцию "posts"
	return &Store{db: db}, nil
}

// Posts возвращает все публикации из базы данных.
func (s *Store) Posts() ([]storage.Post, error) {
	cursor, err := s.db.Find(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to find posts: %w", err)
	}
	defer cursor.Close(context.TODO())

	var posts []storage.Post
	for cursor.Next(context.TODO()) {
		var p storage.Post
		if err := cursor.Decode(&p); err != nil {
			return nil, fmt.Errorf("failed to decode post: %w", err)
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// AddPost добавляет новую публикацию в базу данных.
func (s *Store) AddPost(p storage.Post) error {
	if p.CreatedAt == 0 {
		p.CreatedAt = time.Now().Unix()
	}
	_, err := s.db.InsertOne(context.TODO(), p)
	if err != nil {
		return fmt.Errorf("failed to add post: %w", err)
	}
	return nil
}

// UpdatePost обновляет публикацию в базе данных.
func (s *Store) UpdatePost(p storage.Post) error {
	objectId, err := primitive.ObjectIDFromHex(fmt.Sprint(p.ID))
	if err != nil {
		return fmt.Errorf("failed to convert ID to ObjectID: %w", err)
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"title": p.Title, "content": p.Content, "authorname": p.AuthorName, "createdat": p.CreatedAt}}

	_, err = s.db.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

// DeletePost удаляет публикацию из базы данных.
func (s *Store) DeletePost(p storage.Post) error {

	objectId, err := primitive.ObjectIDFromHex(fmt.Sprint(p.ID))
	if err != nil {
		return fmt.Errorf("failed to convert ID to ObjectID: %w", err)
	}
	filter := bson.M{"_id": objectId}
	_, err = s.db.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}
