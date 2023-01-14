package main

import (
	"birthdayreminder/handlers"
	"birthdayreminder/models"
	"github.com/gin-gonic/gin"
	"log"
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

var (
	allInterests = []models.Interest{
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
	handlers.UpdateAge()
	server := gin.New()
	// creating the user
	server.POST("users", handlers.CreateUser)

	//Getting all the users
	server.GET("users", handlers.GetUsers)

	//Adding friends birthday to the user
	server.POST("/users/:email/birthday", handlers.AddBirthday)

	//Update user friends birthday
	server.PATCH("/users/:email/birthday/:id", handlers.UpdateBirthday)

	err := server.Run()
	if err != nil {
		log.Println("Unable to start server: ", err)
	}

}

// sequential // parallel
// synchronous // asynchronous

// ETL - Extract Transform Load
// running a function on frequent interval
