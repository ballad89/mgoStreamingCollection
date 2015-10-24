##header

**bold**

```
package main

import (
    "fmt"
    "gopkg.in/mgo.v2"
    "mgoStreamingCollection"
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