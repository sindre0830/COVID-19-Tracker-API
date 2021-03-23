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

func (database *Database) Get() error {
	iter := database.Client.Collection("notification").Documents(database.Ctx)
	var notification Notification
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return err
		}
		data := doc.Data()
		//convert data from interface and set in structure
		str := fmt.Sprintf("%v", data["id"])
		notification.ID = str
		str = fmt.Sprintf("%v", data["url"])
		notification.URL = str
		num := data["timeout"].(int)
		notification.Timeout = num
		str = fmt.Sprintf("%v", data["information"])
		notification.Information = str
		str = fmt.Sprintf("%v", data["country"])
		notification.Country = str
		str = fmt.Sprintf("%v", data["trigger"])
		notification.Trigger = str
		//add strucutre to map
		Notifications[notification.ID] = notification
	}
	return nil
}
