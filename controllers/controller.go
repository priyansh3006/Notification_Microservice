package controller

import (
	"context"
	"net/http"
	"notify/database"
	"notify/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var UserCollection *mongo.Collection = database.GetCollection("notifyusers", database.Db)

var validate = validator.New()

func CreateTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var taskInput struct {
			Title    string              `json:"title" validate:"required"`
			Deadline models.TaskDeadline `json:"deadline" validate:"required"`
		}

		if err := c.ShouldBindJSON(&taskInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(taskInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ist, err := time.LoadLocation("Asia/Kolkata")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't load IST location"})
			return
		}

		// Create the deadline in IST
		deadlineIST := taskInput.Deadline.Time().In(ist)
		newTask := models.Task{
			Id:       primitive.NewObjectID(),
			Title:    taskInput.Title,
			Deadline: deadlineIST,
		}

		result, err := UserCollection.InsertOne(ctx, newTask)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't insert new Task in database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": result})
	}
}

func GetAllTasks() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := UserCollection.Find(ctx, bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error in fetching all users"})
		}
		defer result.Close(ctx)
		var tasks []models.Task
		for result.Next(ctx) {
			var temp models.Task
			if err := result.Decode(&temp); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "cant get this user"})
				continue
			}
			tasks = append(tasks, temp)
		}
		c.JSON(http.StatusOK, gin.H{"All task with deadlines": tasks})
	}
}
