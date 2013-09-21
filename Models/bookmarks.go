/* 
The bookmarks package
*/
package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type ( // types for the bookmarks database
	Bookmarks struct {
		//  Id		bson.ObjectId 	`json:"id"	bson:"_id"`
		Key     string     `json:"k"	bson:"k"`
		Created time.Time  `json:"c"	bson:"c"`
		Viewed  time.Time  `json:"v"	bson:"v"`
		List    []Bookmark `json:"l"	bson:"l"`
	}

	Bookmark struct {
		//  Id	bson.ObjectId 	`json:"id"	bson:"_id"`
		URL         string `json:"u"	bson:"u"`
		Description string `json:"d"	bson:"d"`
	}

	BookmarkRepo struct {
		Collection *mgo.Collection
	}
  
  	addRequest struct {
    	list Bookmarks
      	reply chan string // Return the key
    }
  
  	listRetrieve struct {
      key string
      reply chan Bookmarks
  	}
)

// Database Collection info
const  listCollection = "lists"


// Retrieve finds a Bookmarks object in the database
// with the value of the Key field given by key
func retrieve(lists BookmarkRepo, key string) Bookmarks {

  	result := Bookmarks{}
  	err := lists.Collection.Find(bson.M{"key": key}).One(&result)
	if err != nil {
		panic(err)
	}
  
	return result
  
}


// save inserts a Bookmarks object into the database
func save(lists BookmarkRepo, bookmarks Bookmarks) {

  	// Store the current time
    created := time.Now()
  	bookmarks.Created = created
  	bookmarks.Viewed = created
  
  	// Enter the bookmarks document into the database
	err := lists.Collection.Insert(bookmarks)
	if err != nil {
		panic(err)
	}

}

// BookmarksCollection maintains the bookmarks collection
// and serves requests. It provides channels for retrieval,
// insertion, and removal. It communicates to WordList via
// the code channel.
func BookmarksCollection(newColRequest chan ColRequest,
                         getList chan listRetrieve,
                         addList chan addRequest,
                         removeList chan string,
                         code chan string) {
  
  // Get the collection from the database
  reply := make(chan *mgo.Collection)
  newColRequest <- ColRequest{listCollection, reply}
  collection := <- reply
  lists := BookmarkRepo{collection}

  // Serve requests
  // 1. Get list by key
  // 2. Add list
  // 3. Remove list (propagate entries to parents)
  for {
    select {
    case req, ok := <- getList: // Retrieve list by key
      if ok {
        req.reply <- retrieve(lists, req.key) // reply by req.key
      } else { // Caller is dead
      	return // end silently
      }
    case req, ok := <- addList:
      if ok {
      	req.list.Key = <- code        // 1. Generate key
        save(lists, req.list)		// 2. Enter list into database
        req.reply <- req.list.Key   // 3. Return key
      } else { // Caller is dead
      	return // end silently
      }
    case _, ok := <- removeList: // Remove list by key
      if ok {
        // To Delete list by key:
        // 1. Retrieve list
        // 2. Modify list dependencies
        // 3. Remove original from database
        // 4. Free key
        _ = 1
      } else { // Caller is dead
      	return // end silently
      }      
    }
  }
  
}

