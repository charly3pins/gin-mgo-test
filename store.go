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


func create(c *gin.Context) {
	session := s.Copy()
	defer session.Close()

	var user User
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&user)
	if err != nil {
		ErrorWithJSON(c, "Incorrect body", http.StatusBadRequest)
		return
	}
	err = validator.Validate(user)
	if err != nil {
		ErrorWithJSON(c, "Error missing fields", http.StatusBadRequest)
		return
	}

	coll := session.DB("store").C("users")
	now := time.Now()
	user.CreateAt = &now

	err = coll.Insert(user)
	if err != nil {
		if mgo.IsDup(err) {
			ErrorWithJSON(c, "User with this username already exists", http.StatusBadRequest)
			return
		}

		ErrorWithJSON(c, "Database error", http.StatusInternalServerError)
		log.Println("Failed inserting user: ", err)
		return
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusCreated, "User created")
}

func allUsers(c *gin.Context) {
	session := s.Copy()
	defer session.Close()

	coll := session.DB("store").C("users")

	var users []User
	err := coll.Find(bson.M{}).Select(bson.M{"_id": 0, "password": 0, "createat": 0, "updateat": 0}).All(&users)
	if err != nil {
		ErrorWithJSON(c, "Database error", http.StatusInternalServerError)
		log.Println("Failed getting all users: ", err)
		return
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, users)
}

func userByUsername(c *gin.Context) {
	session := s.Copy()
	defer session.Close()

	username := c.Param("username")

	coll := session.DB("store").C("users")

	var user User
	err := coll.Find(bson.M{"username": username}).Select(bson.M{"_id": 0, "password": 0, "createat": 0, "updateat": 0}).One(&user)
	if err != nil || user.Username == "" {
		ErrorWithJSON(c, "User not found", http.StatusNotFound)
		log.Println("Failed finding user: ", err)
		return
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, user)
}

func update(c *gin.Context) {
	session := s.Copy()
	defer session.Close()

	username := c.Param("username")

	var user User
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&user)
	if err != nil {
		ErrorWithJSON(c, "Incorrect body", http.StatusBadRequest)
		return
	}
	err = validator.Validate(user)
	if err != nil {
		ErrorWithJSON(c, "Error missing fields", http.StatusBadRequest)
		return
	}

	coll := session.DB("store").C("users")

	now := time.Now()
	user.UpdateAt = &now

	err = coll.Update(bson.M{"username": username}, &user)
	if err != nil {
		switch err {
		default:
			ErrorWithJSON(c, "Database error", http.StatusInternalServerError)
			log.Println("Failed updating user: ", err)
			return
		case mgo.ErrNotFound:
			ErrorWithJSON(c, "User not found", http.StatusNotFound)
			return
		}
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, user)
}

func delete(c *gin.Context) {
	session := s.Copy()
	defer session.Close()

	username := c.Param("username")

	coll := session.DB("store").C("users")

	err := coll.Remove(bson.M{"username": username})
	if err != nil {
		switch err {
		default:
			ErrorWithJSON(c, "Database error", http.StatusInternalServerError)
			log.Println("Failed deleting username: ", err)
			return
		case mgo.ErrNotFound:
			ErrorWithJSON(c, "User not found", http.StatusNotFound)
			return
		}
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, nil)

}