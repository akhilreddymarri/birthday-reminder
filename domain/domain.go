package domain

import (
	"birthdayreminder/models"
	"errors"
	"fmt"
	"time"
)

var users = make(map[string]models.User)

var (
	ErrUserAlreadyExist     = errors.New("user already exists")
	ErrEmailNotFound        = errors.New("email not found")
	ErrBirthdayAlreadyFound = errors.New("birthday already found")
	ErrBirthdayNotFound     = errors.New("birthday not found")
)

func CreateUser(user models.User) error {
	//check the user is already Exits
	_, exists := isExists(user.Email)
	if exists {
		return ErrUserAlreadyExist
	}
	//creating the id
	user.ID = int64(len(users) + 1)
	users[user.Email] = user

	return nil
}

func GetUsers() map[string]models.User {
	return users
}

func CreateBirthday(email string, birthday models.Birthday) error {
	user, exists := isExists(email)
	if !exists {
		return ErrEmailNotFound
	}
	for _, existingBirthday := range user.FriendsBirthday {
		if existingBirthday.UniqueIdentifier == fmt.Sprintf("%s_%s_%s", birthday.FirstName, birthday.LastName, birthday.Gender) {
			return ErrBirthdayAlreadyFound
		}
	}
	birthday.UniqueIdentifier = fmt.Sprintf("%s_%s_%s", birthday.FirstName, birthday.LastName, birthday.Gender)
	birthday.ID = int64(len(user.FriendsBirthday) + 1)
	user.FriendsBirthday = append(user.FriendsBirthday, birthday)
	users[user.Email] = user
	return nil
}

func UpdateBirthday(id int64, email string, birthday models.Birthday) error {
	// check if given email is present and if not present throw email not present
	user, exists := isExists(email)
	if !exists {
		return ErrEmailNotFound
	}

	isBirthdayExist := false
	for i, existingBirthday := range user.FriendsBirthday {
		if existingBirthday.ID == id {
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
		//context.JSON(http.StatusBadRequest, map[string]string{
		//	"message": "birthday not found",
		//})
		return ErrBirthdayNotFound
	}
	return nil
}

func UpdateAge() {
	for _, user := range users {
		for i, birthday := range user.FriendsBirthday {
			birthday.Age = int(time.Since(birthday.DateOfBirth).Hours() / 24 / 365)
			user.FriendsBirthday[i] = birthday
		}
	}
}

func isExists(email string) (models.User, bool) {
	user, exists := users[email]
	return user, exists
}
