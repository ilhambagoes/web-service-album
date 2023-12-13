package main

import (
	"fmt"
	"net/http"
	"reflect"
	"sort"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// album slices to seed record album data.
var albums = []Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body
func postAlbums(c *gin.Context) {
	var newAlbum Album

	// call BindJSON to bind the received JSON to newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// Parameter sent by the client, then returns that album as a response/
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func deleteAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Declare array response for after delete
	res := make([]Album, 0)

	// convert array albums[] to map
	albumsMap := make(map[int]Album)
	for i := 0; i < len(albums); i++ {
		albumsMap[i+1] = albums[i]
	}

	// Loop over the list of albums, looking for an album whose ID value mathces the parameter.
	for i, keysMap := range albumsMap {
		key := reflect.ValueOf(keysMap)
		idMap := key.FieldByName("ID")

		fmt.Println(id, idMap)

		if id == idMap.String() {
			delete(albumsMap, i)

			for _, key := range albumsMap {
				res = append(res, key)
			}

			sort.Slice(res, func(i, j int) bool {
				return res[i].ID < res[j].ID
			})

			c.IndentedJSON(http.StatusOK, res)
			break
		}
	}
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.DELETE("albums/:id", deleteAlbumByID)

	router.Run("localhost:8080")
}
