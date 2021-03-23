package notification

import (
	"context"
	"fmt"

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

func (database *Database) Add(notification Notification) error {
	_, _, err := database.Client.Collection("notification").Add(database.Ctx, notification)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) Get() (map[string]Notification, error) {
	iter := database.Client.Collection("notification").Documents(database.Ctx)
	var notifications = map[string]Notification {}
	var notification Notification
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, err
		}
		data := doc.Data()
		//convert data from interface and set in structure
		notification.ID = fmt.Sprintf("%v", data["ID"])
		notification.URL = fmt.Sprintf("%v", data["URL"])
		notification.Timeout = data["Timeout"].(int64)
		notification.Information = fmt.Sprintf("%v", data["Information"])
		notification.Country = fmt.Sprintf("%v", data["Country"])
		notification.Trigger = fmt.Sprintf("%v", data["Trigger"])
		//add strucutre to map
		notifications[notification.ID] = notification
	}
	return notifications, nil
}
