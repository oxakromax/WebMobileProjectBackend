package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func clientMongo() *mongo.Client {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://root:52895390@webmobileproject.hcpzpkj.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

type MultipleImages struct {
	UserName string   `json:"username"`
	Images   []string `json:"images"`
}

func main() {
	app := fiber.New()
	app.Post("/upload", UploadRoutine())
	app.Post("/delete", DeleteRoutine())
	app.Get("/getimages", GetRoutine())
	err := app.Listen(":6000")
	if err != nil {
		return
	}
}

func GetRoutine() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// receive a MultipleImages json object and parse it
		var multipleImages MultipleImages
		if err := c.BodyParser(&multipleImages); err != nil {
			return err
		}
		// get the mongo client
		client := clientMongo()
		// get the database
		db := client.Database("WebMobileProject")
		// get the collection
		collection := db.Collection("images")
		// check if the user exists
		var userExists MultipleImages
		err := collection.FindOne(context.TODO(), bson.M{"username": multipleImages.UserName}).Decode(&userExists)
		if err != nil {
			// if the user doesn't exist, return error
			return c.Status(404).JSON(bson.M{"status": "user not found"})
		} else {
			// if the user exists, return "images": userExists.Images, and status: "ok"
			return c.Status(200).JSON(bson.M{
				"status": "ok",
				"images": userExists.Images,
			})
		}
	}
}

func DeleteRoutine() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// receive a MultipleImages json object and parse it
		var multipleImages MultipleImages
		if err := c.BodyParser(&multipleImages); err != nil {
			return err
		}
		// get the mongo client
		client := clientMongo()
		// get the database
		db := client.Database("WebMobileProject")
		// get the collection
		collection := db.Collection("images")
		// check if the user exists
		var userExists MultipleImages
		err := collection.FindOne(context.TODO(), bson.M{"username": multipleImages.UserName}).Decode(&userExists)
		if err != nil {
			// if the user doesn't exist, return error
			return c.Status(404).JSON(bson.M{"status": "user not found"})
		} else {
			// if the user exists, delete the images
			_, err := collection.UpdateOne(context.TODO(), bson.M{"username": multipleImages.UserName}, bson.M{"$pull": bson.M{"images": bson.M{"$in": multipleImages.Images}}})
			if err != nil {
				return err
			}
		}
		// return ok
		return c.Status(200).JSON(bson.M{"status": "ok"})
	}
}

func UploadRoutine() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// receive a MultipleImages json object and parse it
		var multipleImages MultipleImages
		if err := c.BodyParser(&multipleImages); err != nil {
			return err
		}
		// get the mongo client
		client := clientMongo()
		// get the database
		db := client.Database("WebMobileProject")
		// get the collection
		collection := db.Collection("images")
		// check if the user exists
		var userExists MultipleImages
		err := collection.FindOne(context.TODO(), bson.M{"username": multipleImages.UserName}).Decode(&userExists)
		if err != nil {
			// if the user doesn't exist, create it
			_, err := collection.InsertOne(context.TODO(), multipleImages)
			if err != nil {
				return err
			}
		} else {
			// if the user exists, add the new images to the existing ones
			_, err := collection.UpdateOne(context.TODO(), bson.M{"username": multipleImages.UserName}, bson.M{"$push": bson.M{"images": bson.M{"$each": multipleImages.Images}}})
			if err != nil {
				return err
			}
		}
		// return ok
		return c.Status(200).JSON(bson.M{"status": "ok"})
	}
}
