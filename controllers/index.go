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

// func getUser(c *fiber.Ctx) {
// 	collection, err := getMongoDbCollection(dbName, collectionName)
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	var filter bson.M = bson.M{}

// 	if c.Params("id") != "" {
// 		id := c.Params("id")
// 		objID, _ := primitive.ObjectIDFromHex(id)
// 		filter = bson.M{"_id": objID}
// 	}

// 	var results []bson.M
// 	cur, err := collection.Find(context.Background(), filter)
// 	defer cur.Close(context.Background())

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	cur.All(context.Background(), &results)

// 	if results == nil {
// 		c.SendStatus(404)
// 		return
// 	}

// 	json, _ := json.Marshal(results)
// 	c.Send(json)
// }

// func updateUser(c *fiber.Ctx) {
// 	collection, err := getMongoDbCollection(dbName, collectionName)
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	var user User
// 	json.Unmarshal([]byte(c.Body()), &user)

// 	update := bson.M{
// 		"$set": user,
// 	}

// 	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
// 	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	response, _ := json.Marshal(res)
// 	c.Send(response)
// }

// func deleteUser(c *fiber.Ctx) {
// 	collection, err := getMongoDbCollection(dbName, collectionName)

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
// 	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	jsonResponse, _ := json.Marshal(res)
// 	c.Send(jsonResponse)
// }

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

// func createDocument(c *fiber.Ctx) {
// 	collection, err := getMongoDbCollection(dbName, collectionName)
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	var document document.Document
// 	json.Unmarshal([]byte(c.Body()), &document)

// 	res, err := collection.InsertOne(context.Background(), document)
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	response, _ := json.Marshal(res)
// 	c.Send(response)
// }

// func getDocument(c *fiber.Ctx) {
// 	collection, err := getMongoDbCollection(dbName, collectionName)
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	var filter bson.M = bson.M{}

// 	if c.Params("id") != "" {
// 		id := c.Params("id")
// 		objID, _ := primitive.ObjectIDFromHex(id)
// 		filter = bson.M{"_id": objID}
// 	}

// 	var results []bson.M
// 	cur, err := collection.Find(context.Background(), filter)
// 	defer cur.Close(context.Background())

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	cur.All(context.Background(), &results)

// 	if results == nil {
// 		c.SendStatus(404)
// 		return
// 	}

// 	json, _ := json.Marshal(results)
// 	c.Send(json)
// }

// func updateDocument(c *fiber.Ctx) {
// 	collection, err := getMongoDbCollection(dbName, collectionName)
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	var document Document
// 	json.Unmarshal([]byte(c.Body()), &document)

// 	update := bson.M{
// 		"$set": document,
// 	}

// 	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
// 	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	response, _ := json.Marshal(res)
// 	c.Send(response)
// }

// func deleteDocument(c *fiber.Ctx) {
// 	collection, err := getMongoDbCollection(dbName, collectionName)

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
// 	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	jsonResponse, _ := json.Marshal(res)
// 	c.Send(jsonResponse)
// }
