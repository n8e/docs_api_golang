package controllers

import (
	"context"
	"docs_api_golang/configs"
	"docs_api_golang/models"
	"docs_api_golang/responses"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var documentCollection *mongo.Collection = configs.GetCollection(configs.DB, "documents")

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

	// get authenticated user
	var dbUser models.UserSchema
	tokenUser := c.Locals("user").(*jwt.Token)
	claims := tokenUser.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&dbUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "Not authorised to create this document"}})
	}

	newDocument := models.DocumentSchema{
		Id:           primitive.NewObjectID(),
		OwnerId:      dbUser.Id, // should be obtained from authenticated user
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

	err := documentCollection.FindOne(ctx, bson.M{"id": objID}).Decode(&document)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DocumentResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// get authenticated user
	var dbUser models.UserSchema
	tokenUser := c.Locals("user").(*jwt.Token)
	claims := tokenUser.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	readErr := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&dbUser)
	if readErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "Not authorised to create this document"}})
	}

	if document.OwnerId != dbUser.Id {
		return c.Status(http.StatusForbidden).JSON(responses.DocumentResponse{Status: http.StatusForbidden, Message: "error", Data: &fiber.Map{"data": "The logged in user is not authorisd to view this document"}})
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

	// get the document to check owner before updating
	err := documentCollection.FindOne(ctx, bson.M{"id": objID}).Decode(&document)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DocumentResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// get authenticated user
	var dbUser models.UserSchema
	tokenUser := c.Locals("user").(*jwt.Token)
	claims := tokenUser.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	readErr := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&dbUser)
	if readErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "Not authorised to create this document"}})
	}

	if document.OwnerId != dbUser.Id {
		return c.Status(http.StatusForbidden).JSON(responses.DocumentResponse{Status: http.StatusForbidden, Message: "error", Data: &fiber.Map{"data": "The logged in user is not authorisd to update this document"}})
	}

	update := bson.M{"ownerId": document.OwnerId, "title": document.Title, "content": document.Content, "dateCreated": document.DateCreated, "lastModified": document.LastModified}

	result, err := documentCollection.UpdateOne(ctx, bson.M{"id": objID}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DocumentResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// get updated document details
	var updatedDocument models.DocumentSchema
	if result.MatchedCount == 1 {
		err := documentCollection.FindOne(ctx, bson.M{"id": objID}).Decode(&updatedDocument)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.DocumentResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.DocumentResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedDocument}})
}

func DeleteDocument(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	documentId := c.Params("id")
	var document models.DocumentSchema
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(documentId)

	// get the document to check owner before deleting
	err := documentCollection.FindOne(ctx, bson.M{"id": objID}).Decode(&document)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DocumentResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// get authenticated user
	var dbUser models.UserSchema
	tokenUser := c.Locals("user").(*jwt.Token)
	claims := tokenUser.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	readErr := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&dbUser)
	if readErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "Not authorised to create this document"}})
	}

	if document.OwnerId != dbUser.Id {
		return c.Status(http.StatusForbidden).JSON(responses.DocumentResponse{Status: http.StatusForbidden, Message: "error", Data: &fiber.Map{"data": "The logged in user is not authorisd to delete this document"}})
	}

	result, err := documentCollection.DeleteOne(ctx, bson.M{"id": objID})
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

	// get authenticated user
	var dbUser models.UserSchema
	tokenUser := c.Locals("user").(*jwt.Token)
	claims := tokenUser.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	readErr := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&dbUser)
	if readErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "Not authorised to create this document"}})
	}

	results, err := documentCollection.Find(ctx, bson.M{"ownerid": dbUser.Id})
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
