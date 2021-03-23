package firebase

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
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
	defer database.Client.Close()
	return nil
}