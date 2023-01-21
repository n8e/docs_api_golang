package controllers

import (
	"context"
	"docs_api_golang/configs"
	"docs_api_golang/models"
	"docs_api_golang/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var documentCollection *mongo.Collection = configs.GetCollection(configs.DB, "documents")
var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.UserSchema
	defer cancel()

	// validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newUser := models.UserSchema{
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("id")
	var user models.UserSchema
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})
}

func UpdateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("id")
	var user models.UserSchema
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(userId)

	// validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// use the validator library to validate reqiured fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"userName": user.UserName, "firstName": user.FirstName, "lastName": user.LastName, "email": user.Email, "password": user.Password, "role": user.Role}

	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// get updated user details
	var updatedUser models.UserSchema
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedUser)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
}

func DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("id")
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "User with specified ID not found!"}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}})
}

func GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.UserSchema
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.UserSchema
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}})
}

// DOCUMENTS

func CreateDocument(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var document models.DocumentSchema
	defer cancel()

	// validate the request body
	if err := c.BodyParser(&document); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.DocumentResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// use the validator library to validate required fields
	if validationErr := validate.Struct(&document); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.DocumentResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newDocument := models.DocumentSchema{
		OwnerId:      document.OwnerId,
		Title:        document.Title,
		Content:      document.Content,
		DateCreated:  document.DateCreated,
		LastModified: document.LastModified,
	}

	result, err := documentCollection.InsertOne(ctx, newDocument)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DocumentResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.DocumentResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetDocument(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	documentId := c.Params("id")
	var document models.DocumentSchema
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(documentId)

	err := documentCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&document)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DocumentResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.DocumentResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": document}})
}

func UpdateDocument(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	documentId := c.Params("id")
	var document models.DocumentSchema
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(documentId)

	// validate the request body
	if err := c.BodyParser(&document); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.DocumentResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// use the validator library to validate reqiured fields
	if validationErr := validate.Struct(&document); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.DocumentResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"ownerId": document.OwnerId, "title": document.Title, "content": document.Content, "dateCreated": document.DateCreated, "lastModified": document.LastModified}

	result, err := documentCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DocumentResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// get updated document details
	var updatedDocument models.DocumentSchema
	if result.MatchedCount == 1 {
		err := documentCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedDocument)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.DocumentResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.DocumentResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedDocument}})
}

func DeleteDocument(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	documentId := c.Params("id")
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(documentId)

	result, err := documentCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DocumentResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusInternalServerError).JSON(responses.DocumentResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "Document with specified ID not found!"}})
	}

	return c.Status(http.StatusOK).JSON(responses.DocumentResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "Document successfully deleted!"}})
}

func GetDocuments(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var documents []models.DocumentSchema
	defer cancel()

	results, err := documentCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DocumentResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleDocument models.DocumentSchema
		if err = results.Decode(&singleDocument); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.DocumentResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		documents = append(documents, singleDocument)
	}

	return c.Status(http.StatusOK).JSON(responses.DocumentResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": documents}})
}
