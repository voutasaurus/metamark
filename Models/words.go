/* 
The words package provides the means to initialise and maintain
a collection of words, and query the collection for a unique short 
string of words.

The package also provides a means of ensuring that the strings are
unique, via UniqueCodeTracker.

In your main function:

Initialise:

request := make(chan ColRequest)
quit := make(chan bool)
go Database(request chan ColRequest, quit chan bool)
defer close(request)
defer close(quit)

word := make(chan string) // For getting a random word
go WordList(randWord, request)
defer close(randWord)

newCode := make(chan string) // For getting a new unique code
freeCode := make(chan string) // For freeing a code after deletion
go UniqueCodeTracker(newCode, freeCode, randWord)
defer close(newCode)
defer close(freeCode)

During execution:

When you need to get a random word:
word <- randWord

When you need to get a new code:
code <- newCode

When you need to free a code:
freeCode <- code

*/

package models

import (
	//"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
  	"math/rand"
)

type (
  Word struct {
  	Key	int		`json:"k"	bson:"k"`
    Word string `json:"w"	bson:"w"`
  }

  WordRepo struct {
  	Collection *mgo.Collection
  }
  
  CodeRequest struct {
  	reply chan string
  }
)

// Database Collection info
const  wordCollection = "word"

// WordList initialises and maintains the goside of
// the words collection from the database. It servesstring
// requests for random words via the word channel.
// These requests come from UniqueCodeTracker.
func WordList(word chan string, newColRequest chan ColRequest) {

  // Get the collection from the database
  reply := make(chan *mgo.Collection)
  newColRequest <- ColRequest{wordCollection, reply}
  collection := <- reply
  words := WordRepo{collection}
  
  // Find the number of words
  max, err := words.Collection.Count()
  if err != nil {
  	panic(err)
  }
  
  r := rand.New(rand.NewSource(time.Now().UnixNano())) // random generator
  result := Word{} // result of word query
  
  for { // Serve forever
    // Get a random word in anticipation
    err = words.Collection.Find(bson.M{"key": r.Intn(max)}).One(&result)
    if err != nil {
      panic(err)
    }
    select { // Wait for either a word request or a closed channel
    case word <- result.Word:
    case _, ok := <- word:
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
func UniqueCodeTracker(newCode chan string, freeCode chan string, word chan string) {
  
  var codesInUse = make(map[string]bool)
  codesInUse[""] = true // required for the generate function
  
  var unique string
  needNewCode := true
  
  for {
    // Pregenerate a unique code
    if (needNewCode) { // Only generate if last unique code has been served
	    codesInUse, unique = generate(codesInUse, word)
      	needNewCode = false
    }
    select {
    case newCode <- unique:
      	needNewCode = true
    case code, ok := <- freeCode:
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
  n := 8 // number of words in unique code
  var str string // stores unique code
 
  for codesInUse[str] { // If the code is in use, try again
    
    for i := 0; i < n; i++ {
      str += <- word // Get a word and add it to the string
      if i < n - 1 {
          str += " "
      }
 	} // A candidate code has been generated

  } // A unique code has been generated
  
  // Mark code as in use
  codesInUse[str] = true
  
  return codesInUse, str
  
}
