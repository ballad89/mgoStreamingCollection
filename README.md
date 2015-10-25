##MongoDB streaming collection
This package is built ontop of the mongodb driver `gopkg.in/mgo.v2`.
It adds functionality to
- Check if a collection exist
- Convert an uncapped collection to a capped one
- A constructor like function, which returns a capped collection
- Get a collection's stats
- Query a capped collection and stream documents out of it into a channel 

**Installation**
`go get github.com/ballad89/mgoStreamingCollection`

**Example Usage**

Examples of how to use the individual functions can be seen in the test file
Below is an example of how you would use it in a program

```
package main

import (
    "fmt"
    "gopkg.in/mgo.v2"
    "github.com/ballad89/mgoStreamingCollection"
)

func main() {
    connectUrl := "mongodb://33.33.13.45/streamingDb"

    session, err := mgo.Dial(connectUrl)

    if err != nil {
        panic(err)
    }

    session.SetSafe(&mgo.Safe{})

    defer session.Close()

    collectionName := "cappedCollection"
    collectionSize := 1024

    database := session.DB("streamingDB")

    cappedCollection, err := mgoStreamingCollection.CreateOrConvertCollection(database, collectionName, collectionSize)

    if err != nil {
        panic(err)
    }

    ch := make(chan interface{})

    go mgoStreamingCollection.TailQuery(cappedCollection.Find(nil), ch)

    for msg := range ch {
        fmt.Println(msg)
    }
}
```