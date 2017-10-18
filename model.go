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


type User struct {
	ID       bson.ObjectId `json:"ID,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name" validate:"nonzero"`
	Username string        `json:"username" validate:"nonzero"`
	Password string        `json:"password,omitempty" validate:"nonzero"`
	Roles    []int         `json:"roles,omitempty"`
	Email    string        `json:"email" validate:"nonzero"`
	Language string        `json:"language" validate:"nonzero"`
	CreateAt *time.Time    `json:"createat,omitempty"`
	UpdateAt *time.Time    `json:"updateat,omitempty"`
}