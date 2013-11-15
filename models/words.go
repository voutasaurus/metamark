// Copyright comment

package models

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"math/rand"
	"time"
)

type (
	Word struct {
		//Key  int    "k"
		Word string "w"
	}
/*
	WordRepo struct {
		Collection *mgo.Collection
	}
*/
	CodeRequest struct {
		reply chan string
	}
)

// Database Collection info
const wordCollection = "word"
const wordField = "w"

// getwords reads the words from the database into a local variable stored on the heap.
// This cuts down on database accesses, and encapsulates the database interactions, 
// making the random word selection more portable.
func getwords(localwords []string, words *mgo.Collection) error {
	// Within the words collection, Find documents with the word field, and then Select only
  	// the contents of that field, and return an Iterator.
  	iter := words.Find(bson.M{wordField: bson.M{"$exists": true}}).Select(bson.M{wordField: 1}).Limit(10000).Iter()
  
  	// Iterate through the results pulling out the strings.
  tempresult := Word{}	
	for i:=0; iter.Next(&tempresult); i++ {
        localwords[i] = tempresult.Word
      //Debug: fmt.Print(tempresult.Word, ",")
    }
  
  //Debug: fmt.Println(len(localwords))
  	// We should check len(localwords) to catch errors
  	if len(localwords) < 1000 {
  		// some error - we need to close the iterator also.
      	// composite errors may occur.
  	}
  	// Close the iterator, return an error where applicable.
    if err := iter.Close(); err != nil {
        return err
    }
  	return nil	
}

// WordList initialises and maintains the goside of
// the words collection from the database. It servesstring
// requests for random words via the word channel.
// These requests come from UniqueCodeTracker.
func WordList(word chan string, newColRequest chan ColRequest, quit chan bool) {
  
	// Get the collection from the database
	reply := make(chan *mgo.Collection)
	newColRequest <- ColRequest{wordCollection, reply}
	words := <-reply // this is the collection

  	// Read the entire word collection into a slice.
  	// Approx 10000 words, average 8 runes each, rune = 32 bits = 4 bytes,
  	// total approx 32KB.
  // TODO: Should make this global mutexed(*) for all WordList goroutines to use (at the moment there is only
  	// one WordList gorountine).
  	// Multiple reads can be simultaneous, but simultaneous read & write should be disallowed.
  	localwords := make([]string, 10000)
  	err := getwords(localwords, words)
  	if err != nil {
  		panic(err) // Replace this with error handling function
  	}
  
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // random generator
 
	for { // Serve forever
      
		// Get a random word in anticipation.
      	// Note that localwords may be refreshed inbetween uses and len is constant.
		v := r.Intn(len(localwords)) 	// random integer in [0,len(localwords))
      	tempword := localwords[v] 		// Random word in localwords

		//Debug
		fmt.Println(tempword)      
      
      	// Wait for either a word request or a quit request
		select { 
		case word <- tempword:
        case <-quit:
			return // silently end
		}
      	// TODO: add this: case: refresh localwords
	}

}

// UniqueCodeTracker initialises and maintains the register
// of which codes are in use. It serves requests for new unique
// codes, and requests to free expired codes. word is a
// channel for communicating with WordList.
func UniqueCodeTracker(newCode chan string, freeCode chan string, colRequest chan ColRequest) {

	// Run the words collection
	randWord := make(chan string) // For getting a random word
  	quit := make(chan bool)
	go WordList(randWord, colRequest, quit)
	defer close(randWord)
  	defer close(quit)

	// Keep track of which codes are in use
  // TODO: Make this global to all UniqueCodeTracker goroutines (presently there is only one UCT goroutine)
	var codesInUse = make(map[string]bool)
	codesInUse[""] = true // required for the generate function

	// Unique code
	var unique string
	needNewCode := true

	// Serve requests
	// First pregenerate a unique code
	// 1. Send the generated code when receiver is ready
	// 2. Remove the code provided to freeCode
	for {
		// Pregenerate a unique code
		if needNewCode { // Only generate if last unique code has been served
			codesInUse, unique = generate(codesInUse, randWord)
			needNewCode = false
		}
		select {
		case newCode <- unique: // 1. Send new code
			needNewCode = true
		case code, ok := <-freeCode: // 2. Delete expired code
			if ok { // Caller wants to free a code
				delete(codesInUse, code)
			} else { // Caller is dead
              	// tell wordlist to stop
              	quit <- true
				return
			}
		}
	} // loop back forever

}

// generate is a helper for UniqueCodeTracker. It generates a unique
// code by communicating with WordList.
func generate(codesInUse map[string]bool, word chan string) (map[string]bool, string) {

	//TODO: make n = 6, 7, 8 randomly with appropriate weightings
	//		given the distribution of max^6, max^7 and max^8
	n := 8         // number of words in unique code
	var str string // stores unique code

	for codesInUse[str] { // If the code is in use, try again

		for i := 0; i < n; i++ {
			x := <-word
			str += x // Get a word and add it to the string
			if i < n-1 {
				str += " "
			}
		} // A candidate code has been generated

	} // A unique code has been generated

	// Mark code as in use
	codesInUse[str] = true

	return codesInUse, str

}
