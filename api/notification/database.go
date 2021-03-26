package notification

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// DB stores information about the firebase database.
//
// Used in notification structure.
var DB Database

// Database structure stores information about firebase database.
//
// Functionality: Setup, Add, Get, Delete
type Database struct {
	Ctx context.Context
	Client *firestore.Client
}

// Setup sets up database.
func (database *Database) Setup() error {
	//connect to firebase with key and branch if an error occurred
	database.Ctx = context.Background()
	sa := option.WithCredentialsFile("./serviceAccountKey.json")
	app, err := firebase.NewApp(database.Ctx, nil, sa)
	if err != nil {
		return err
	}
	//open firestore and branch if an error occurred
	database.Client, err = app.Firestore(database.Ctx)
	if err != nil {
		return err
	}
	//update Notifications with data from database and branch if and error occurred
	err = database.Get()
	if err != nil {
		return err
	}
	return nil
}

// Add adds a new webhook to database.
func (database *Database) Add(notification *Notification) error {
	//add webhook to database and get a UUID from firebase and branch if an error occurred
	key, _, err := database.Client.Collection("notification").Add(database.Ctx, notification)
	if err != nil {
		return err
	}
	//update webhook ID in database with UUID and branch if an error occurred
	notification.ID = key.ID
	_, err = database.Client.Collection("notification").Doc(key.ID).Update(database.Ctx, []firestore.Update {{
		Path:  "ID",
		Value: notification.ID,
	}})
	if err != nil {
		return err
	}
	//update Notifications with data from database and branch if and error occurred
	err = database.Get()
	if err != nil {
		return err
	}
	return nil
}

// Get gets all webhooks from database
func (database *Database) Get() error {
	//clear current webhooks stored in Notifications
	Notifications = make(map[string]Notification)
	//iterate through database and get each webhook
	iter := database.Client.Collection("notification").Documents(database.Ctx)
	var notification Notification
	for {
		//go to next element in array and break loop if there are no elements, branch if an error occurred
		elem, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return err
		}
		//convert data from interface and set in structure
		data := elem.Data()
		notification.ID = fmt.Sprintf("%v", data["ID"])
		notification.URL = fmt.Sprintf("%v", data["URL"])
		notification.Timeout = data["Timeout"].(int)
		notification.Information = fmt.Sprintf("%v", data["Information"])
		notification.Country = fmt.Sprintf("%v", data["Country"])
		notification.Trigger = fmt.Sprintf("%v", data["Trigger"])
		//add structure to map
		Notifications[notification.ID] = notification
	}
	return nil
}

// Delete deletes specific webhook.
func (database *Database) Delete(id string) error {
	//get only element that has the same ID as specified and branch if an error occurred
	iter := database.Client.Collection("notification").Where("ID", "==", id).Documents(database.Ctx)
	elem, err := iter.Next()
	if err != nil {
		return err
	}
	//delete webhook and branch if an error occurred
	_, err = elem.Ref.Delete(database.Ctx)
	if err != nil {
		return err
	}
	//update Notifications with data from database and branch if and error occurred
	err = database.Get()
	if err != nil {
		return err
	}
	return nil
}
