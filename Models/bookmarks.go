// Copyright

package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

// types for the bookmarks database
type (

	// Bookmarks stores a list of bookmarks.
	// The Key which is generated when the object is passed to the collection.
	// The Created and Viewed fields are maintained by the collection.
	// The user of the package should initialise the List field.
	Bookmarks struct {
		//  Id		bson.ObjectId 	`json:"id"	bson:"_id"`
		Key     string     `json:"k"	bson:"k"`
		Created time.Time  `json:"c"	bson:"c"`
		Viewed  time.Time  `json:"v"	bson:"v"`
		List    []Bookmark `json:"l"	bson:"l"`
	}

	// Bookmark stores a single bookmark.
	// The user of the package should initialise the URL and Description fields.
	Bookmark struct {
		//  Id	bson.ObjectId 	`json:"id"	bson:"_id"`
		URL         string `json:"u"	bson:"u"`
		Description string `json:"d"	bson:"d"`
	}

	// bookmarkRepo stores the collection
	bookmarkRepo struct {
		Collection *mgo.Collection
	}

	// AddRequest stores a request for adding a Bookmarks object
	// to the bookmarks collection
	AddRequest struct {
		List  Bookmarks
		Reply chan string // Return the key
	}

	// ListRetrieve stores a request for retrieving a Bookmarks object
	// from the bookmarks collection by key
	ListRetrieve struct {
		Key   string
		Reply chan Bookmarks // Returns the Bookmarks object found
	}
)

// Database Collection info
const listCollection = "lists"

// Retrieve finds a Bookmarks object in the database
// with the value of the Key field given by key
func retrieve(lists bookmarkRepo, key string) Bookmarks {

	result := Bookmarks{}
	err := lists.Collection.Find(bson.M{"k": key}).One(&result)
	if err != mgo.ErrNotFound {
		if err != nil {
			panic(err)
		}
	}

	return result

}

// save inserts a Bookmarks object into the database
func save(lists bookmarkRepo, bookmarks Bookmarks) {

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
func BookmarksCollection(getList chan ListRetrieve,
	addList chan AddRequest,
	removeList chan string) {
	// Run database
	newColRequest := make(chan ColRequest)
	dbQuit := make(chan bool)
	go Database(newColRequest, dbQuit)
	defer close(newColRequest)
	defer close(dbQuit)

	// Run the code generator - opens the words collection
	newCode := make(chan string)  // For getting a new unique code
	freeCode := make(chan string) // For freeing a code after deletion
	go UniqueCodeTracker(newCode, freeCode, newColRequest)
	defer close(newCode)
	defer close(freeCode)

	// Get the bookmarks collection from the database
	reply := make(chan *mgo.Collection)
	newColRequest <- ColRequest{listCollection, reply}
	collection := <-reply
	lists := bookmarkRepo{collection}

	// Serve requests
	// 1. Get list by key
	// 2. Add list
	// 3. Remove list (propagate entries to parents)
	for {
		select {
		case req, ok := <-getList: // Retrieve list by key
			if ok {
				req.Reply <- retrieve(lists, req.Key) // reply by req.key
			} else { // Caller is dead
				return // end silently
			}
		case req, ok := <-addList:
			if ok {
				req.List.Key = <-newCode  // Generate key
				save(lists, req.List)     // Enter list into database
				req.Reply <- req.List.Key // Return key
			} else { // Caller is dead
				return // end silently
			}
		case _, ok := <-removeList: // Remove list by key
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
