package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

// Globals
var (
	isHeroku      = checkHeroku()
	configuration = loadConfig()
	db            = initDB()
)

var roomConnectionMap = make(map[string][]*websocket.Conn)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Connection success!")
}

func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	// generate the room string
	roomString := randString(4)

	// insert the new room
	queryString := "INSERT INTO rooms(room_code, start_time) VALUES($1, now())"
	stmt, err := db.Prepare(queryString)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusInternalServerError), w, http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(roomString)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusInternalServerError), w, http.StatusInternalServerError)
		return
	}

	// generate a secret to share with the creator
	roomSecret := randString(32)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, roomSecret+","+roomString)
}

func joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	// unmarhall the room string
	decoder := json.NewDecoder(r.Body)
	req := struct {
		RoomString string
	}{""}

	err := decoder.Decode(&req)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		return
	}

	// upgrade to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		return
	}

	message := struct{ stringAck string }{""}
	// frontend handshake to get user and hook them into the userMap for sockets
	err = conn.ReadJSON(&message)
	if err != nil {
		failWithStatusCode(err, "Failed to handshake", w, http.StatusInternalServerError)
		return
	}

	roomConnectionMap[req.RoomString] = append(roomConnectionMap[req.RoomString], conn)

	w.WriteHeader(http.StatusOK)
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	// unmarhall the question id
	decoder := json.NewDecoder(r.Body)
	req := struct {
		RoomID     string
		QuestionID string
	}{"", ""}

	err := decoder.Decode(&req)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
	}

	queryString := "UPDATE questions SET votes = votes + 1 WHERE q_id = $1"
	stmt, err := db.Prepare(queryString)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusInternalServerError), w, http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(req.QuestionID)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusInternalServerError), w, http.StatusInternalServerError)
		return
	}

	// update them sockets
	for _, socket := range roomConnectionMap[req.RoomID] {
		fmt.Println(socket)
		// grab all the questions from the database for this room and send them back over the socket
	}

	w.WriteHeader(http.StatusOK)
}

func hideHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func loadConfig() Configuration {
	configuration := Configuration{}
	if !isHeroku {
		file, err := os.Open("conf.json")
		failOnError(err, "Config json not found. Make sure it is present.")
		decoder := json.NewDecoder(file)

		err = decoder.Decode(&configuration)
		if err != nil {
			fmt.Println("error:", err)
		}
	}
	return configuration
}

func initDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", configuration.DB.Host, configuration.DB.Port, configuration.DB.User, configuration.DB.Pass, configuration.DB.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	failOnError(err, "Failed to open Postgres")

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(i) * time.Second)

		if err = db.Ping(); err == nil {
			break
		}
		log.Println(err)
	}

	if err != nil {
		failGracefully(err, "Failed to open Postgres")
	}
	err = db.Ping()
	if err != nil {
		failGracefully(err, "Failed to Ping Postgres")
	} else {
		fmt.Println("Connected to DB")
	}
	return db
}

func main() {
	fmt.Printf("Listening on port: %s\n", configuration.Port)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/createRoom", createRoomHandler)
	http.HandleFunc("/joinRoom", joinRoomHandler)
	http.HandleFunc("/vote", voteHandler)
	http.HandleFunc("/hide", hideHandler)
	http.ListenAndServe(configuration.Port, nil)
}
