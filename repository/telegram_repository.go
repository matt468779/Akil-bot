package repository

import (
	"akil_telegram_bot/domain"
	"akil_telegram_bot/mongo"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type telegramRepository struct {
	database   mongo.Database
	collection string
}

func NewTelegramRepository(db mongo.Database, collection string) domain.TelegramRepository {
	return &telegramRepository{
		database:   db,
		collection: collection,
	}
}

func (tr *telegramRepository) SaveMessage(c context.Context, update *domain.Update) error {
	coll := tr.database.Collection(tr.collection)
	_, err := coll.InsertOne(c, update)

	return err
}

func (tr *telegramRepository) GetMessages(c context.Context, update *domain.Update) []domain.Update {
	coll := tr.database.Collection(tr.collection)
	filter := bson.M{
		"message.chat.id": update.Message.Chat.Id,
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "message.date", Value: -1}})
	findOptions.SetLimit(10)
	cursor, err := coll.Find(c, filter, findOptions)
	if err != nil {
		log.Println(err)
	}
	var result []domain.Update

	for cursor.Next(context.Background()) {
		var resUpdate domain.Update
		err := cursor.Decode(&resUpdate)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, resUpdate)
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}
