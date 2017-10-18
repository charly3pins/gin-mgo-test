package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"

	"github.com/gin-gonic/gin"
)

var s *mgo.Session

func main() {
	var err error
	s, err = mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer s.Close()

	s.SetMode(mgo.Monotonic, true)
	ensureIndex(s)

	router := gin.Default()
	router.POST("/signup", create)
	router.GET("/users", allUsers)
	router.GET("/users/:username", userByUsername)
	router.PUT("/users/:username", update)
	router.POST("/users/:username", delete)

	router.Run(":8080")
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("store").C("users")

	index := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}
