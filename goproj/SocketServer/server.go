package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Latitude    float64            `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude   float64            `json:"longitude,omitempty" bson:"longitude,omitempty"`
	GroundSpeed float64            `json:"groundspeed" bson:"groundspeed"`
	UTCTime     time.Time          `json:"utctime,omitempty" bson:"utctime,omitempty"`
}

const (
	message       = "ping"
	StopCharacter = "\r\n\r\n"
)

func SocketServer(port int) {

	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))

	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", port, err)
		os.Exit(1)
	}

	defer listen.Close()

	log.Printf("Begin listen port: %d", port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go handler(conn)
	}

}

func handler(conn net.Conn) {

	defer conn.Close()

	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
		w   = bufio.NewWriter(conn)
	)

ILOOP:
	for {
		n, err := r.Read(buf)
		data := string(buf[:n])

		switch err {
		case io.EOF:
			break ILOOP
		case nil:
			log.Println("Receive:", data)
			if isTransportOver(data) {
				break ILOOP
			}

		default:
			log.Fatalf("Receive data failed:%s", err)
			return
		}

		// Declare host and port options to pass to the Connect() method
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		fmt.Println("clientoptions type:", reflect.TypeOf(clientOptions), "\n")
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			fmt.Println("mongo.Connect() ERROR:", err)
			os.Exit(1)
		}
		ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

		// Access a MongoDB collection through a database
		col := client.Database("golang").Collection("socket_data")
		fmt.Println("Collection type:", reflect.TypeOf(col), "\n")

		// Declare a MongoDB struct instance for the document's fields and data

		// Should print out 'main.MongoFields' struct

		var decodedata Post
		error := json.Unmarshal(buf[:n], &decodedata)
		if error != nil {
			fmt.Println("json convert ERROR:", error)
			os.Exit(1)
		}
		fmt.Println("oneDoc TYPE:", reflect.TypeOf(decodedata), "\n")
		result, insertErr := col.InsertOne(ctx, decodedata)
		if insertErr != nil {
			fmt.Println("InsertOne ERROR:", insertErr)
			os.Exit(1) // safely exit script on error
		} else {

			fmt.Println("InsertOne() result type: ", reflect.TypeOf(result))
			fmt.Println("InsertOne() API result:", result)

			// get the inserted ID string
			newID := result.InsertedID
			fmt.Println("InsertOne() newID:", newID)
			fmt.Println("InsertOne() newID type:", reflect.TypeOf(newID))
		}

	}
	w.Write([]byte(message))
	w.Flush()
	log.Printf("Send: %s", message)

}

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\r\n\r\n")
	return
}
func deleteinvaliddata() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	fmt.Println("clientoptions type:", reflect.TypeOf(clientOptions), "\n")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		os.Exit(1)
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// Access a MongoDB collection through a database
	col := client.Database("golang").Collection("socket_data")
	fmt.Println("Collection type:", reflect.TypeOf(col), "\n")
	findOptions := options.Find()
	cur, err := col.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		fmt.Println("cursor:", cur)
		var curdoc Post
		if err = cur.Decode(&curdoc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(curdoc.UTCTime)
		t := time.Now()
		before30 := t.Add(-30 * 24 * time.Hour)
		if curdoc.UTCTime.Before(before30) {
			result, error := col.DeleteOne(ctx, bson.M{"utctime": curdoc.UTCTime})
			if error != nil {
				fmt.Println("DeleteDocment Error:", error)
			}
			fmt.Println("Deleted documents number:", result.DeletedCount)
		}

	}
	// create a value into which the single document can be decoded

}

func main() {

	port := 3333
	deleteinvaliddata()
	SocketServer(port)

}
