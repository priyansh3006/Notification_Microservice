package main

import (
	"context"
	"fmt"
	"log"
	"time"

	controller "notify/controllers"
	"notify/models"
	"notify/notification"
	"notify/routes"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection = controller.UserCollection
var router *gin.Engine

var loc, _ = time.LoadLocation("Asia/Kolkata")

func retrieveAllData() {
	log.Println("Starting retrieveAllData")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error querying all documents: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		log.Printf("Error decoding documents: %v", err)
		return
	}

	log.Printf("Found %d documents", len(tasks))
	for _, task := range tasks {
		log.Printf("Task: %+v", task)
	}
	log.Println("Finished retrieveAllData")
}

func checkDeadlinesAndNotify() {
	log.Println("Starting checkDeadlinesAndNotify")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now().In(loc)
	log.Printf("Current time: %v", now)

	// Delete tasks with passed deadlines
	deleteResult, err := collection.DeleteMany(ctx, bson.M{
		"task_deadline": bson.M{"$lt": now},
	})
	if err != nil {
		log.Printf("Error deleting past tasks: %v", err)
	} else {
		log.Printf("Deleted %d tasks with passed deadlines", deleteResult.DeletedCount)
	}

	// Find tasks with upcoming deadlines
	cursor, err := collection.Find(ctx, bson.M{
		"task_deadline": bson.M{"$gt": now},
	})
	if err != nil {
		log.Printf("Error querying tasks: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		log.Printf("Error decoding tasks: %v", err)
		return
	}

	log.Printf("Found %d tasks", len(tasks))

	for _, task := range tasks {
		log.Printf("Checking task: %s, Deadline: %v", task.Title, task.Deadline)
		if task.ShouldNotify() {
			log.Printf("Should notify for task: %s", task.Title)
			message := fmt.Sprintf("Reminder: Task '%s' is due in less than 1 hour", task.Title)
			deletePassedTask, err := collection.DeleteOne(ctx, bson.M{"task_title": task.Title})
			if err != nil {
				log.Printf("Error deleting task %s from database: %v as it have already been notified", task.Title, err)
				continue
			}
			fmt.Println(deletePassedTask.DeletedCount)
			if err := notification.EmailSender("priyansh3006@example.com", "Task Reminder", message); err != nil {
				log.Printf("Error sending email for task %s: %v", task.Title, err)
			} else {
				log.Printf("Email sent for task: %s", task.Title)
			}
			fmt.Println(message)
			if err := notification.SendSms("7023045653", message); err != nil {
				log.Printf("Error sending SMS for task %s: %v", task.Title, err)
			} else {
				log.Printf("SMS sent for task: %s", task.Title)
			}
		}
	}
	log.Println("Finished checkDeadlinesAndNotify")
}

func main() {
	var port = ":8080"
	router = gin.New()
	router.Use(gin.Logger())
	routes.SetupRoutes(router)

	fmt.Println("Starting the server")
	log.Println("Starting the server")

	c := cron.New()
	c.AddFunc("* * * * *", checkDeadlinesAndNotify)
	log.Println("Cron job added")
	c.Start()
	log.Println("Cron job started")

	retrieveAllData()

	// Run the server in a separate goroutine
	go func() {
		fmt.Println("Inside the go routine now")
		log.Println("Inside the go routine now") // Added log statement
		if err := router.Run(port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Main function completed, entering select{}")
	// Keep the application running
	select {}
}
