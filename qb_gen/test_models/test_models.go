package test_models

import (
	"time"

	"github.com/wyattis/goof/gtime"
)

type User struct {
	Id   int
	Name string
}

type Comment struct {
	Id        int
	UserId    int
	Body      string
	CreatedAt time.Time
	UpdatedAt gtime.TimeRFC1123
}

type Activity struct {
	Id       int
	SecretId string
}

func (a Activity) TableName() string {
	return "activities"
}

func (a *Activity) SetDefaultSecretId() error {
	return nil
}
