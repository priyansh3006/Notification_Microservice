package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var loc, _ = time.LoadLocation("Asia/Kolkata")

type TaskDeadline struct {
	Year   int `json:"year" bson:"year" validate:"required"`
	Month  int `json:"month" bson:"month" validate:"required,min=1,max=12"`
	Day    int `json:"day" bson:"day" validate:"required,min=1,max=31"`
	Hour   int `json:"hour" bson:"hour" validate:"required,min=0,max=23"`
	Minute int `json:"minute" bson:"minute" validate:"required,min=0,max=59"`
}

func (td TaskDeadline) Time() time.Time {
	return time.Date(td.Year, time.Month(td.Month), td.Day, td.Hour, td.Minute, 0, 0, time.Local)
}

type Task struct {
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"task_title" validate:"required"`
	Deadline time.Time          `json:"deadline" bson:"task_deadline" validate:"required"`
}

func (t *Task) ShouldNotify() bool {
	deadlineTime := t.Deadline
	now := time.Now().In(loc)
	timeDiff := deadlineTime.Sub(now)
	return timeDiff > 0 && timeDiff <= time.Hour
}
