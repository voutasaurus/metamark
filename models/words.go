// Copyright comment

package models

import (
	"fmt"
	"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"
	"math/rand"
	"time"
)

type (
	Word struct {
		Key  int    "k"
		Word string "w"
	}

	WordRepo struct {
		Collection *mgo.Collection
	}

	CodeRequest struct {
		reply chan string
	}
)

// Database Collection info
const wordCollection = "word"

// WordList initialises and maintains the goside of
// the words collection from the database. It servesstring
// requests for random words via the word channel.
// These requests come from UniqueCodeTracker.
func WordList(word chan string, newColRequest chan ColRequest) {

	// Get the collection from the database
	reply := make(chan *mgo.Collection)
	newColRequest <- ColRequest{wordCollection, reply}
	collection := <-reply
	words := WordRepo{collection}

	r := rand.New(rand.NewSource(time.Now().UnixNano())) // random generator
	result := Word{}                                     // result of word query

	for { // Serve forever
      
      // Find the number of words - this may change between requests
      max, err := words.Collection.Count()
  	  // fmt.Println(max)
      if err != nil {	// TODO: catalogue credible errors and handle them appropriately
          panic(err)
      }
      
      	// Get a random word in anticipation
		v := r.Intn(max)
		err = words.Collection.Find(bson.M{"k": v}).One(&result)
		//err = words.Collection.Find(nil).Skip(v).One(&result)
		if err != nil {
          // TODO: Replace with switch statement - maybe in a separate function
          if err != mgo.ErrNotFound { 
				panic(err)
			} else {
				fmt.Println("not found")
			}
		}
		fmt.Println(result.Word)
		select { // Wait for either a word request or a closed channel
		case word <- result.Word:
		case _, ok := <-word:
			if !ok { // caller is dead
				return // silently end
			}
		}
	}

}

// UniqueCodeTracker initialises and maintains the register
// of which codes are in use. It serves requests for new unique
// codes, and requests to free expired codes. word is a
// channel for communicating with WordList.
func UniqueCodeTracker(newCode chan string, freeCode chan string, colRequest chan ColRequest) {

	// Run the words collection
	randWord := make(chan string) // For getting a random word
	go WordList(randWord, colRequest)
	defer close(randWord)

	// Keep track of which codes are in use
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
