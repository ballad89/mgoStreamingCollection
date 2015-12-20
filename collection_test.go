package mgoStreamingCollection_test

import (
	"github.com/ballad89/mgoStreamingCollection"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestConnection(t *testing.T) {

	database, err := setup()

	if err != nil {
		t.Fatalf("error connecting to database with error: %v", err)

	}

	defer database.Session.Close()

	collectionName := "testConnection"
	collectionSize := 1024

	_, err = mgoStreamingCollection.CreateOrConvertCollection(database, collectionName, collectionSize)

	if err != nil {
		t.Errorf("collection creation failed with error: %v", err)
	}
}

func TestCollectionExistsAndStats(t *testing.T) {
	database, err := setup()

	if err != nil {
		t.Fatalf("error connecting to database with error: %v", err)

	}

	defer database.Session.Close()

	testColl := "testStats"
	colSize := 8192

	coll, err := mgoStreamingCollection.CreateOrConvertCollection(database, testColl, colSize)

	if err != nil {
		t.Errorf("error creating collection with error: %v", err)
	}

	exists, err := mgoStreamingCollection.CollectionExists(database, testColl)

	if err != nil {
		t.Errorf("error checking collection existence with error: %v", err)
	}

	if !exists {
		t.Errorf("collection %s does not exist", testColl)
	}

	stats, err := mgoStreamingCollection.Stats(coll)

	if err != nil {
		t.Errorf("error getting collection stats with error: %v", err)
	}

	if !stats.Capped {
		t.Errorf("collection %s not capped", testColl)
	}

	if stats.MaxBytes != colSize {
		t.Log(stats.MaxBytes)
		t.Errorf("collection %s size not set %d", testColl, colSize)
	}

}

func TestConvertCollection(t *testing.T) {

	database, err := setup()

	if err != nil {
		t.Fatalf("error connecting to database with error: %v", err)

	}

	defer database.Session.Close()

	testColl := "testConvert"
	colSize := 8192

	coll := database.C(testColl)
	coll.DropCollection()

	coll.Insert(bson.M{"name": "husayn"})

	stats, err := mgoStreamingCollection.Stats(coll)

	if err != nil {
		t.Errorf("error getting collection stats with error: %v", err)
	}

	if stats.Capped {
		t.Errorf("collection %s should not be capped", testColl)
	}

	err = mgoStreamingCollection.ConvertToCapped(coll, colSize)

	if err != nil {
		t.Errorf("error converting collection with error: %v", err)
	}

	newStats, err := mgoStreamingCollection.Stats(coll)

	if err != nil {
		t.Errorf("error getting collection stats with error: %v", err)
	}

	if !newStats.Capped {
		t.Errorf("collection %s should be capped", testColl)
	}

	if newStats.MaxBytes != colSize {
		t.Errorf("collection %s size not set %d", testColl, colSize)
	}

}

func TestTailQuery(t *testing.T) {

	database, err := setup()

	if err != nil {
		t.Fatalf("error connecting to database with error: %v", err)

	}

	defer database.Session.Close()

	collectionName := "testTailQuery"
	collectionSize := 1024

	if exists, _ := mgoStreamingCollection.CollectionExists(database, collectionName); exists {
		database.C(collectionName).DropCollection()
	}

	cappedCollection, err := mgoStreamingCollection.CreateOrConvertCollection(database, collectionName, collectionSize)

	if err != nil {
		t.Errorf("collection creation failed with error: %v", err)
	}

	ch := make(chan interface{})

	go mgoStreamingCollection.TailQuery(cappedCollection.Find(nil), ch)

	for i := 0; i < 3; i++ {
		cappedCollection.Insert(bson.M{"name": "husayn", "count": i})
		for msg := range ch {
			if msg != nil {
				t.Log(msg)
				break
			}
			t.Errorf("inserted document not read with cound %d", i)
		}
	}
	close(ch)

}

func setup() (*mgo.Database, error) {

	connectUrl := "mongodb://33.33.13.45/test"

	session, err := mgo.Dial(connectUrl)

	if err != nil {
		return nil, err

	}

	session.SetSafe(&mgo.Safe{})

	database := session.DB("test")

	return database, nil

}
