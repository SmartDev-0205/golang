package main

import (
	"encoding/json"
	// manage multiple requests
	"errors" // os.Exit(1) on Error
	// get an object type

	"fmt" // Println() function

	"log"
	"net"
	"os"

	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrNoRecord           = errors.New("models:no matching record found")
	ErrInvalidCredentials = errors.New("models:invalis credentials")
	ErrDuplicateEmail     = errors.New("models:duplicate email")
)

type Post struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Latitude    float64            `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude   float64            `json:"longitude,omitempty" bson:"longitude,omitempty"`
	GroundSpeed float64            `json:"groundspeed," bson:"groundspeed,"`
	UTCTime     time.Time          `json:"utctime,omitempty" bson:"utctime,omitempty"`
}

const (
	StopCharacter = "\r\n\r\n"
)

func SocketClient(ip string, port int) {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	defer conn.Close()
	Message := &Post{
		Latitude:    39.000679,
		Longitude:   121.400269,
		GroundSpeed: 0.000000,
		UTCTime:     time.Now(),
	}
	MsgStr, error := json.Marshal(Message)
	if error != nil {
		fmt.Println(error)
		return
	}
	conn.Write([]byte(MsgStr))
	conn.Write([]byte(StopCharacter))
	log.Printf("Send: %s", MsgStr)

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])

}

func main() {
	var (
		ip   = "127.0.0.1"
		port = 3333
	)
	SocketClient(ip, port)

}
