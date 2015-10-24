package mgoStreamingCollection_test

import (
	"gopkg.in/mgo.v2"
	"mgoStreamingCollection"
	"testing"
)

func TestConnection(t *testing.T) {
	connectUrl := "mongodb://33.33.13.45/streamingDb"

	session, err := mgo.Dial(connectUrl)

	if err != nil {
		t.Fatalf("error connecting to database with error: %v", err)
	}

	session.SetSafe(&mgo.Safe{})

	defer session.Close()

	collectionName := "cappedCollection9"
	collectionSize := 1024

	database := session.DB("streamingDB")

	cappedCollection, err := mgoStreamingCollection.CreateOrConvertCollection(database, collectionName, collectionSize)

	if err != nil {
		t.Errorf("collection creation failed with error: %v", err)
	}

	mgoStreamingCollection.Stats(cappedCollection)
}
