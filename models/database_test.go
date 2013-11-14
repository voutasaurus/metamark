
package models

import "testing"

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