package models

import (
	"fmt"
	//"html/template"
	//"io/ioutil"
	//"net/http"
	//"regexp"
	//"strings"
	//"errors"
  	"labix.org/v2/mgo"
  	"labix.org/v2/mgo/bson"
  	"time"
)

type (
  Bookmarks struct {
  //  Id		bson.ObjectId 	`json:"id"	bson:"_id"`
    Key		string			`json:"k"	bson:"k"`
    Created	time.Time		`json:"c"	bson:"c"`
    Viewed	time.Time		`json:"v"	bson:"v"`
   	List 	[]Bookmark 		`json:"l"	bson:"l"`
  }
  
  Bookmark struct {
  //  Id	bson.ObjectId 	`json:"id"	bson:"_id"`
    URL	string 			`json:"u"	bson:"u"`
    Description	string 	`json:"d"	bson:"d"`
  }
  
  BookmarkRepo struct {
  	Collection *mgo.Collection
  }
)

const (
	blue = "blueDB"
	collection = "lists"
)


var repo BookmarkRepo{ session.DB(blue).C(collection) }

func save(bookmarks Bookmarks) {

  err = repo.Collection.Insert(bookmarks)
  if err != nil {
  	panic(err)
  }
  
}


func init() {
  	url := "192.168.1.10:27017"
  	session, err := mgo.Dial(url) 
}

func test() {
  
    db := "modelsTest"
    col := "initTest"
    c := session.DB(db).C(col)
  
  	
  err = c.Insert(&Bookmarks{Key: "testkey", Created: time.Now(), Viewed: time.Now()})
  if err != nil {
    panic(err)
  }
  
 	result := Bookmarks{}
  key := "testkey"
  err = c.Find(bson.M{"key": key}).One(&result)
  if err != nil {
    panic(err)
  }

  fmt.Println(result.Created)

}

/*

func (r BookmarkRepo) All() (bookmarks Bookmarks, err error) {
    err = r.Collection.Find(bson.M{}).All(&bookmarks)
    return
}

func handleBookmarks(w http.ResponseWriter, r *http.Request) {
    var (
        bookmarks Bookmarks
        err   error
    )
    if bookmarks, err = repo.All(); err != nil {
        log.Printf("%v", err)
        http.Error(w, "500 Internal Server Error", 500)
        return
    }
    writeJson(w, bookmarks)
}
*/
