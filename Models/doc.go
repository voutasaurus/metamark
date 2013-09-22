/*
The models package provides a means of communicating with the database

In your main function Initialise:

getList := make(chan ListRetrieve)
addList := make(chan AddRequest)
removeList := make(chan string)
go BookmarksCollection(getList, addList, removeList)

This runs the bookmarks collection and provides channels for requests.

To make requests of the database, do the following:

1. Retrieve a list

key := keyToRetrieve
reply := make(chan Bookmarks)
getList <- ListRetrieve{key, reply}
newList := <- reply

2. Add a list

bookmarks := bookmarksToInsert
reply := make(chan string)
addList <- AddRequest{bookmarks, reply}
newKey := <- reply

3. Remove a list

key := keyToDelete
removeList <- key

*/
package models
