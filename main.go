package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

//users
//	birthdays
//		name
//		dob
//		gender
//		profile pic
//		interests
//			books
//			sports
//			video games

type User struct {
	ID              int64      `json:"ID"`
	FirstName       string     `json:"firstName"`
	LastName        string     `json:"lastName"`
	FriendsBirthday []Birthday `json:"friendsBirthday"`
}

type Birthday struct {
	ID         int64     `json:"ID"`
	Name       string    `json:"name"`
	DOB        time.Time `json:"DOB"`
	Gender     string    `json:"gender"`
	ProfilePic string    `json:"profilePic"`
	Interests  []int64   `json:"interests"`
}

type Interest struct {
	ID       int64  `json:"ID"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

var (
	allInterests = []Interest{
		{
			ID:       1,
			Name:     "Books",
			Category: "Art",
		},
		{
			Name:     "Painting",
			Category: "Art",
		},
		{
			Name:     "Sports",
			Category: "Art",
		},
	}
)

func main() {
	server := gin.New()
	err := server.Run()
	if err != nil {
		log.Println("Unable to start server: ", err)
	}
}
