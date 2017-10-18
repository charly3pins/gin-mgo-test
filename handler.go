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

func ErrorWithJSON(c *gin.Context, message string, code int) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(code, message)
	fmt.Printf("{message: %q}", message)
}
