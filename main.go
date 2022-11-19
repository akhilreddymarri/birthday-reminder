package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
	Email           string     `json:"email"`
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

var users = make(map[string]User)

func main() {
	server := gin.New()
	server.POST("user", func(context *gin.Context) {
		var user User
		//binding the json
		err := context.ShouldBindJSON(&user)
		if err != nil {
			log.Println("Error while creating the user")
			context.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
			return
		}

		//check the user is already Exits
		user, isExits := users[user.Email]
		if isExits {
			context.JSON(http.StatusConflict, map[string]string{
				"message": "user already exists",
			})
			return
		}
		//creating the id
		user.ID = int64(len(users) + 1)
		users[user.Email] = user
		context.JSON(http.StatusOK, map[string]string{
			"message": "success",
		})
	})
	err := server.Run()
	if err != nil {
		log.Println("Unable to start server: ", err)
	}
}
