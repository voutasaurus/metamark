package models

import (
	"fmt"
	"html/template"
	//"io/ioutil"
	"net/http"
	"regexp"
	//"strings"
	//"errors"
  	"labix.org/v2/mgo"
  	"time"
)

type (
  Bookmarks struct {
    Id		bson.ObjectId 	`json:"id"	bson:"_id"`
    Created	time.Time		`json:"c"	bson:"c"`
    Viewed	time.Time		`json:"v"	bson:"v"`
   	List 	[]Bookmark 		`json:"c"	bson:"c"`
  }
  
  Bookmark struct {
    Id	bson.ObjectId 	`json:"id"	bson:"_id"`
    URL	string 			`json:"u"	bson:"u"`
    Description	string 	`json:"d"	bson:"d"`
  }
)