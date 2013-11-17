// Copyright

package models

import (
	//"fmt"
	"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"
)

const (
	url = "192.168.1.2:27017" // mongodb address
	//url  = "localhost:27017" // temp mongodb address for testing
	blue = "blueDB"
)

// Request for a handle of a collection
type ColRequest struct {
	col   string               // collection name
	reply chan *mgo.Collection // handle for session connection to the collection
}

/*

construct colReq
go Database(colReq)
defer(close(colReq))

*/

// Database starts a session on the database and provides channels
// for other functions to take control over certain collections
func Database(newColRequest chan ColRequest) {

	session, err := mgo.Dial(url) // Connect to the database
	if err != nil {               // If database is not active
		panic(err)
		// How should one recover from database failure?
		// Maybe this should call mongod somehow to reestablish the database?
	}
	defer session.Close() // closes the session when database is returned

	for {
		select {
		case req, ok := <-newColRequest:
			if ok { // Caller wants a new collection
				req.reply <- session.DB(blue).C(req.col)
			} else { // Caller is dead - channel is closed
				return // silently end
			}
		}
	}

}
