package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"math/big"
	"math/rand"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

const (
	REDIS_KEY = "UNIQUE_NUMBER_MAP"
	PORT      = 8080
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// global context
var ctx = context.Background()

// init redis client
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379", // Redis server address
	DB:   0,                // Use default DB
})

// generateUniqueBigInt: generate random integer then check that whether it exists in redis or not, then update a map.
func generateUniqueBigInt() *big.Int {
	for {
		num := new(big.Int).Rand(rand.New(rand.NewSource(time.Now().UnixNano())), big.NewInt(1e18))
		numStr := num.String()

		exists, err := rdb.SIsMember(ctx, REDIS_KEY, numStr).Result()
		if err != nil {
			log.Println("Redis error:", err)
			continue
		}

		if !exists {
			// add number to redis
			err := rdb.SAdd(ctx, REDIS_KEY, numStr).Err()
			if err != nil {
				log.Println("Redis error:", err)
				continue
			}
			return num
		}
	}
}

// handleConnection: handle connection in web socket
func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", message)

		uniqueBigInt := generateUniqueBigInt()
		err = conn.WriteMessage(websocket.TextMessage, []byte(uniqueBigInt.String()))
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)

	log.Printf("Server started on %d", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
