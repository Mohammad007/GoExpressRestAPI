package database

import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/Mohammad007/GoExpressRestAPI/internal/models"
)

type MongoDB struct {
    client *mongo.Client
    db     *mongo.Database
}

func NewMongoDB(config Config) (*MongoDB, error) {
    uri := fmt.Sprintf("mongodb://%s:%s", config.Host, config.Port)
    client, err := mongo.NewClient(options.Client().ApplyURI(uri))
    if err != nil {
        return nil, err
    }
    return &MongoDB{client: client, db: client.Database(config.DBName)}, nil
}

func (m *MongoDB) Connect() error {
    return m.client.Connect(context.Background())
}

func (m *MongoDB) Close() error {
    return m.client.Disconnect(context.Background())
}

func (m *MongoDB) CreateUser(ctx context.Context, user *models.User) error {
    collection := m.db.Collection("users")
    _, err := collection.InsertOne(ctx, user.UserSchema)
    return err
}

func (m *MongoDB) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
    collection := m.db.Collection("users")
    var user models.User
    err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&user.UserSchema)
    return &user, err
}

func (m *MongoDB) GetAllUsers(ctx context.Context) ([]models.User, error) {
    collection := m.db.Collection("users")
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    var users []models.User
    for cursor.Next(ctx) {
        var user models.User
        if err := cursor.Decode(&user.UserSchema); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}

func (m *MongoDB) UpdateUser(ctx context.Context, user *models.User) error {
    collection := m.db.Collection("users")
    _, err := collection.UpdateOne(ctx, bson.M{"id": user.ID}, bson.M{"$set": user.UserSchema})
    return err
}

func (m *MongoDB) DeleteUser(ctx context.Context, id uint) error {
    collection := m.db.Collection("users")
    _, err := collection.DeleteOne(ctx, bson.M{"id": id})
    return err
}