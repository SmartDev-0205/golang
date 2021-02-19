package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goproj/config/db"
	"goproj/model"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	//"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	//jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type Post struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Latitude    float32            `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude   float32            `json:"longitude,omitempty" bson:"longitude,omitempty"`
	GroundSpeed float32            `json:"groundspeed," bson:"groundspeed,"`
	UTCTime     time.Time          `json:"utctime,omitempty" bson:"utctime,omitempty"`
}

var (
	ErrNoRecord           = errors.New("models:no matching record found")
	ErrInvalidCredentials = errors.New("models:invalid credentials")
	ErrDuplicateEmail     = errors.New("models:duplicate email")
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json")
	var user model.User

	body, err := ioutil.ReadAll(r.Body)
	//req := r.URL.Query()
	err = json.Unmarshal(body, &user)
	//mapstructure.Decode(req, &user)
	user.Created = time.Now()
	fmt.Println("param:", user)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return

	}

	collection, err := db.GetUserDBCollection()
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var result model.User
	err = collection.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			hash, err := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), 5)

			if err != nil {
				res.Error = "Error while Hashing Password,Try Again"
				json.NewEncoder(w).Encode(res)
				return
			}

			user.HashedPassword = string(hash)
			_, err = collection.InsertOne(context.TODO(), user)
			if err != nil {
				res.Error = "Error While Crating User,Try Again"
				json.NewEncoder(w).Encode(res)
				return
			}
			res.Result = "Registration Successfully"
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	} else {
		res.Error = ErrDuplicateEmail.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application-json")
	var user model.User
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &user)

	fmt.Println("param:", user)

	if err != nil {
		log.Fatal(err)

	}

	collection, err := db.GetUserDBCollection()
	if err != nil {
		log.Fatal(err)
	}
	var result model.User
	var res model.ResponseResult
	err = collection.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&result)
	if err != nil {
		res.Error = ErrNoRecord.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.HashedPassword), []byte(user.HashedPassword))
	if err != nil {
		res.Error = ErrInvalidCredentials.Error()

		json.NewEncoder(w).Encode(res)
		return
	}
	res.Result = "successfully logined"
	res.User = result
	findOptions := options.Find()
	collection, err = db.GetSocketDBCollection()
	if err != nil {
		res.Error = "cannot read SocketDB"

		json.NewEncoder(w).Encode(res)
		return
	}
	fmt.Println("now starting")
	cursor, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("now entering")
	for cursor.Next(context.TODO()) {
		fmt.Println("cursor:", cursor)
		var curdoc model.Post
		if err = cursor.Decode(&curdoc); err != nil {
			fmt.Println("error occured!")
			log.Fatal(err)
		}
		res.MapPosInfo = append(res.MapPosInfo, curdoc)
		fmt.Println("info", res.MapPosInfo)
	}
	json.NewEncoder(w).Encode(res)
	return

}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home started:")
	w.Header().Set("Content-type", "application-json")
	findOptions := options.Find()
	collection, err := db.GetSocketDBCollection()
	var res model.ResponseResult
	if err != nil {
		res.Error = "cannot read SocketDB"

		json.NewEncoder(w).Encode(res)
		return
	}
	fmt.Println("now starting")
	cursor, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("now entering")
	for cursor.Next(context.TODO()) {
		fmt.Println("cursor:", cursor)
		var curdoc model.Post
		if err = cursor.Decode(&curdoc); err != nil {
			fmt.Println("error occured!")
			log.Fatal(err)
		}
		res.MapPosInfo = append(res.MapPosInfo, curdoc)
		fmt.Println("info", res.MapPosInfo)
	}
	json.NewEncoder(w).Encode(res)
	return
}
