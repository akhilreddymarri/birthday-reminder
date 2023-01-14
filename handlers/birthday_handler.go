package handlers

import (
	"birthdayreminder/domain"
	"birthdayreminder/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"time"
)

var fieldValidator *validator.Validate

func init() {
	fieldValidator = validator.New()

	fieldValidator.RegisterValidation("isValidDOB", func(fl validator.FieldLevel) bool {
		dob := fl.Field().Interface().(time.Time)
		currentTime := time.Now()
		return dob.Before(currentTime)
	})
}

func CreateUser(context *gin.Context) {
	var user models.User
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

	err = domain.CreateUser(user)
	if err != nil {
		context.JSON(http.StatusConflict, map[string]string{
			"message": err.Error(),
		})
	}

	context.JSON(http.StatusOK, map[string]string{
		"message": "success",
	})
}

func GetUsers(context *gin.Context) {
	context.JSON(http.StatusOK, map[string]interface{}{
		"users": domain.GetUsers(),
	})
}

func AddBirthday(context *gin.Context) {
	//bind the uri with a golang structure
	var emailParam models.EmailParam
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
	var birthday models.Birthday
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

	err = domain.CreateBirthday(emailParam.Email, birthday)
	if err != nil {
		if err == domain.ErrEmailNotFound {
			context.JSON(http.StatusNotFound, map[string]string{
				"message": err.Error(),
			})
			return
		}
		if err == domain.ErrBirthdayAlreadyFound {
			context.JSON(http.StatusConflict, map[string]string{
				"message": err.Error(),
			})
			return
		}
		context.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, map[string]string{
		"message": "Friends birthday added successfully",
	})

}

func UpdateBirthday(context *gin.Context) {
	// should bind uri with email and id
	var updateParam models.UpdateParam
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
	var birthday models.Birthday
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

	err = domain.UpdateBirthday(updateParam.ID, updateParam.Email, birthday)
	if err != nil {
		if err == domain.ErrEmailNotFound {
			context.JSON(http.StatusNotFound, map[string]string{
				"message": err.Error(),
			})
			return
		}
		if err == domain.ErrBirthdayNotFound {
			context.JSON(http.StatusNotFound, map[string]string{
				"message": err.Error(),
			})
			return
		}
		context.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, map[string]string{
		"message": "updated successfully",
	})
}

func UpdateAge() {
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			select {
			case <-ticker.C:
				domain.UpdateAge()
				log.Println("updating birthday")
			}
		}
	}()
}
