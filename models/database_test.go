package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"testing"
)

const keyfield = "key"
const namefield = "name"

// TestCollectionResponse tests connections with existing collections.
// This includes Create, Read, Update, and Destroy operations.
func TestCollectionResponse(t *testing.T) {
	// Make a colReq channel
	colReq := make(chan ColRequest)

	// Activate database
	go Database(colReq)
	close(colReq)

	// Connect to existing collection
	// Build connection request
	// Send connection request

	// The following tests should be in separate functions:
	//	(Not named TestXXX - they would then be run independently which would be bad)

	// Test Read:
	// Positive: To test Read we try to read something we know exists in
	// the collection and check that it matches.
	// Negative: Also try to read for a document that doesn't exist.

	// Test Create:
	// Positive: To test create, generate a random document, check that it is
	// unique and then put it in the database. Read to see that it is
	// there and matches the input.

	// Test Update:
	// Positive: To test update, pick a document, update it in a random
	// way, read it by id and then compare to see that the changes
	// are correct.
	// Negative: Also try to update an non existent document.

	// Test Destroy:
	// Positive: To test destroy, pick a document to destroy, and destroy it.
	// Read to make sure that it is not there.
	// Negative: Also try to destroy a non existent document.

}

// Object for database
type Sample struct {
	key  int
	name string
}

func Read(col *mgo.Collection, key int) Sample {
	// result holder
	iter := col.Find(bson.M{keyfield: key}).Select(bson.M{keyfield: 1, namefield: 1}).Iter()

	result := Sample{}
	iter.Next(&result)
	iter.Close()

	return result
}

func Create(col *mgo.Collection, s Sample) error {
	return col.Insert(s)
}

func Update(col *mgo.Collection, key int, newname string) error {
	return col.Update(bson.M{keyfield: key}, bson.M{namefield: newname})
}

func Destroy(col *mgo.Collection, key int) error {
	return col.Remove(bson.M{keyfield: key})
}
