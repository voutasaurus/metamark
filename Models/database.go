// Copyright

package models

import (
	//"fmt"
	"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"
)

const (
	url  = "192.168.1.10:27017" // mongodb address
	blue = "blueDB"
)

type ColRequest struct {
	col   string
	reply chan *mgo.Collection
}

// Database starts a session on the database and provides channels
// for other functions to take control over certain collections
func Database(newColRequest chan ColRequest, q chan bool) {

	session, err := mgo.Dial(url) // Connect to the database
	if err != nil {
		panic(err)
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
		case <-q: // quit signal
			return // silently end
		}
	}

}
