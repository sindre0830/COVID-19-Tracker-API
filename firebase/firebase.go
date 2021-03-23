package firebase

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var DB Database

type Database struct {
	Ctx context.Context
	Client *firestore.Client
}

func (database *Database) Setup() error {
	// Use a service account
	database.Ctx = context.Background()
	sa := option.WithCredentialsFile("./serviceAccountKey.json")
	app, err := firebase.NewApp(database.Ctx, nil, sa)
	if err != nil {
		return err
	}
	database.Client, err = app.Firestore(database.Ctx)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) Add(notification interface{}) error {
	_, _, err := database.Client.Collection("notification").Add(database.Ctx, notification)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) Get() ([]map[string]interface{}, error) {
	iter := database.Client.Collection("notification").Documents(database.Ctx)
	var output []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, err
		}
		data := doc.Data()
		output = append(output, data)
	}
	return output, nil
}
