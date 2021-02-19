package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	//ID             primitive.ObjectID `json:"_id,omitempty"`
	Name           string    `json:"name,omitempty" `
	Surname        string    `json:"surname,omitempty" `
	Email          string    `json:"email,omitempty" `
	HashedPassword string    `json:"hashedPassword,omitempty" `
	Created        time.Time `json:"created,omitempty" `
	token          string    `json:"token,omitempty" `
}
type Post struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Latitude    float64            `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude   float64            `json:"longitude,omitempty" bson:"longitude,omitempty"`
	GroundSpeed float64            `json:"groundspeed," bson:"groundspeed,"`
	UTCTime     time.Time          `json:"utctime,omitempty" bson:"utctime,omitempty"`
}
type ResponseResult struct {
	Error      string `json:"error"`
	User       User   `json:"user"`
	MapPosInfo []Post `json:"mapposinfo"`
	Result     string `json:"result"`
}
