package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	ID              int64      `json:"id"`
	FirstName       string     `json:"firstName" validate:"required,alpha,min=3,max=50"`
	LastName        string     `json:"lastName" validate:"required,alpha,min=1,max=50"`
	Email           string     `json:"email" validate:"required,email"`
	FriendsBirthday []Birthday `json:"friendsBirthday"`
}

type Birthday struct {
	UniqueIdentifier string    `json:"uniqueIdentifier"`
	ID               int64     `json:"id"`
	FirstName        string    `json:"firstName" validate:"required,alpha,min=3,max=50"`
	LastName         string    `json:"lastName" validate:"required,alpha,min=1,max=50"`
	DateOfBirth      time.Time `json:"dateOfBirth" validate:"isValidDOB"`
	Age              int       `json:"age"`
	Gender           string    `json:"gender" validate:"required,oneof=MALE FEMALE"`
	ProfilePic       string    `json:"profilePic"`
	Interests        []int64   `json:"interests"`
}

type Interest struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type EmailParam struct {
	Email string `uri:"email" validate:"required,email"`
}
type UpdateParam struct {
	ID    int64  `uri:"id" validate:"required"`
	Email string `uri:"email" validate:"required,email"`
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

var fieldValidator *validator.Validate

func init() {
	fieldValidator = validator.New()

	fieldValidator.RegisterValidation("isValidDOB", func(fl validator.FieldLevel) bool {
		t := fl.Field().Interface().(time.Time)
		return t.Before(time.Now())
	})
}

func main() {
	server := gin.New()
	// creating the user
	server.POST("users", func(context *gin.Context) {
		var user User
		//binding the json (converting postman json body to golang struct)
		err := context.ShouldBindJSON(&user)
		if err != nil {
			log.Println("Error while creating the user")
			context.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
			return
		}
		// manual validations performmed below
		//if user.FirstName == "" || user.LastName == "" || user.Email == "" {
		//	context.JSON(http.StatusBadRequest, map[string]string{
		//		"message": "required fields are missing",
		//	})
		//	return
		//}

		err = fieldValidator.Struct(user)
		if err != nil {
			context.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
			return
		}

		//check the user is already Exits
		_, isExits := users[user.Email]
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

	//Getting all the users
	server.GET("users", func(context *gin.Context) {

		context.JSON(http.StatusOK, map[string]interface{}{
			"users": users,
		})
	})

	//Adding friends birthday to the user
	server.POST("/users/:email/birthday", func(context *gin.Context) {
		//bind the uri with a golang structure
		var emailParam EmailParam
		err := context.ShouldBindUri(&emailParam)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		//do field level validations for URI structure (fieldValidator.Struct(email))
		err = fieldValidator.Struct(emailParam)
		if err != nil {
			log.Println("Error while validating the email", err)
			context.JSON(http.StatusBadRequest, map[string]string{
				"message": "Email required",
			})
			return
		}

		//bind the request body(birthday) with a golang structure
		var birthday Birthday
		err = context.ShouldBindJSON(&birthday)
		if err != nil {
			log.Println("Error while adding the birthday")
			context.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
			return
		}

		//do field level validations for request body structure
		err = fieldValidator.Struct(birthday)
		if err != nil {
			log.Println("Error while validating the email", err)
			context.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
			return
		}

		user, isExits := users[emailParam.Email]
		if !isExits {
			context.JSON(http.StatusBadRequest, map[string]string{
				"message": "Email not found",
			})
			return
		}
		for _, existingBirthday := range user.FriendsBirthday {
			if existingBirthday.UniqueIdentifier == fmt.Sprintf("%s_%s_%s", birthday.FirstName, birthday.LastName, birthday.Gender) {
				context.JSON(http.StatusConflict, map[string]string{
					"message": "birthday already found",
				})
				return
			}
		}
		birthday.UniqueIdentifier = fmt.Sprintf("%s_%s_%s", birthday.FirstName, birthday.LastName, birthday.Gender)
		birthday.ID = int64(len(user.FriendsBirthday) + 1)
		user.FriendsBirthday = append(user.FriendsBirthday, birthday)
		users[user.Email] = user

		context.JSON(http.StatusOK, map[string]string{
			"message": "Friends birthday added successfully",
		})

	})

	//Update user friends birthday
	server.PATCH("/users/:email/birthday/:id", func(context *gin.Context) {
		// should bind uri with email and id
		var updateParam UpdateParam
		err := context.ShouldBindUri(&updateParam)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		// do field level validations for URI structure
		err = fieldValidator.Struct(updateParam)
		if err != nil {
			log.Println("Error while validating the id and email", err)
			context.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
			return
		}

		// should bind body for birthday body
		var birthday Birthday
		err = context.ShouldBindJSON(&birthday)
		if err != nil {
			log.Println("Error while updating the birthday", err)
			context.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
			return
		}

		// do field level validations for body structure
		err = fieldValidator.Struct(birthday)
		if err != nil {
			log.Println("Error while validating the birthday", err)
			context.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
			return
		}

		// check if given email is present and if not present throw email not present
		user, isExits := users[updateParam.Email]
		if !isExits {
			context.JSON(http.StatusBadRequest, map[string]string{
				"message": "Email not present",
			})
			return
		}

		isBirthdayExist := false
		for i, existingBirthday := range user.FriendsBirthday {
			if existingBirthday.ID == updateParam.ID {
				isBirthdayExist = true
				user.FriendsBirthday[i].FirstName = birthday.FirstName
				user.FriendsBirthday[i].LastName = birthday.LastName
				user.FriendsBirthday[i].Gender = birthday.Gender
				user.FriendsBirthday[i].DateOfBirth = birthday.DateOfBirth
				user.FriendsBirthday[i].UniqueIdentifier = fmt.Sprintf("%s_%s_%s", birthday.FirstName, birthday.LastName, birthday.Gender)
				users[user.Email] = user
				break
			}
		}

		if !isBirthdayExist {
			context.JSON(http.StatusBadRequest, map[string]string{
				"message": "birthday not found",
			})
			return
		}

		context.JSON(http.StatusBadRequest, map[string]string{
			"message": "updated successfully",
		})
	})

	err := server.Run()
	if err != nil {
		log.Println("Unable to start server: ", err)
	}

}
