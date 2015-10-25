package mgoStreamingCollection

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func CollectionExists(db *mgo.Database, name string) (bool, error) {

	result, err := db.CollectionNames()

	if err != nil {
		return false, err
	}

	for _, v := range result {
		if v == name {
			return true, nil
		}
	}

	return false, nil
}

func ConvertToCapped(c *mgo.Collection, size int) error {
	return c.Database.Run(bson.D{{"convertToCapped", c.Name}, {"size", size}}, nil)
}

func createCappedCollection(db *mgo.Database, name string, size int) error {
	return db.C(name).Create(&mgo.CollectionInfo{
		Capped:   true,
		MaxBytes: size,
	})
}

// Creates a capped collection called `collectionName`.
// if `collectionName` exists but is not capped, it is converted to a capped collection
//
func CreateOrConvertCollection(database *mgo.Database, collectionName string, size int) (*mgo.Collection, error) {

	exists, err := CollectionExists(database, collectionName)

	if err != nil {
		return nil, err
	}

	if !exists {
		err := createCappedCollection(database, collectionName, size)
		if err != nil {
			return nil, err
		}
	}

	collection := database.C(collectionName)

	collectionStats, err := Stats(collection)

	if err != nil {
		return nil, err
	}

	if !collectionStats.Capped || collectionStats.MaxBytes != size {
		err = ConvertToCapped(collection, size)
		if err != nil {
			return nil, err
		}
	}

	return collection, nil
}

func Stats(c *mgo.Collection) (*mgo.CollectionInfo, error) {
	stats := struct {
		Capped  bool
		Maxsize int `bson:"maxSize"`
	}{}

	err := c.Database.Run(bson.D{{"collStats", c.Name}}, &stats)

	if err != nil {
		return nil, err
	}

	result := mgo.CollectionInfo{
		Capped:   stats.Capped,
		MaxBytes: stats.Maxsize,
	}
	return &result, nil
}

// Tails a query on a capped collection
// streams the returned documents over the channel passed in
// This method should be run in a goroutine
// closing the channel ends the goroutine
func TailQuery(query *mgo.Query, ch chan interface{}) {
	defer func() {
		recover()
	}()

	iter := query.Tail(-1)
	defer iter.Close()

	for {
		var result interface{}

		for iter.Next(&result) {
			ch <- result
		}

		if err := iter.Err(); err != nil {
			iter.Close()
		}

		if iter.Timeout() {
			continue
		}

		iter = query.Tail(-1)
	}

}
