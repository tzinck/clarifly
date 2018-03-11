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
	swears        = loadProfanity("en")
)

var roomConnectionMap = make(map[string][]*websocket.Conn)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	// generate the room string
	roomString := randString(4)
	// generate a secret to share with the creator
	roomSecret := randString(32)

	// insert the new room
	queryString := "INSERT INTO rooms(room_code, secret, start_time) VALUES($1, $2, now())"
	stmt, err := db.Prepare(queryString)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusInternalServerError), w, http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(roomString, roomSecret)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusInternalServerError), w, http.StatusInternalServerError)
		return
	}

	fmt.Println("Creating room with code " + roomString + ".")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, roomSecret+","+roomString)
}

func joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	// upgrade to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		return
	}

	messageType, p, err := conn.ReadMessage()
	if err != nil {
		failWithStatusCode(err, "Failed to handshake", w, http.StatusInternalServerError)
		return
	}

	// Check DB if room exists
	var code string
	err = db.QueryRow("SELECT room_code FROM rooms WHERE room_code = $1", string(p)).Scan(&code)
	if err == sql.ErrNoRows || err != nil {
		// Room does not exist
		returnmsg := []byte("Room " + string(p) + " does not exist.")
		err = conn.WriteMessage(messageType, returnmsg)

		if err != nil {
			fmt.Println("error message broke bad lol")
			return
		}

		fmt.Println("Room " + string(p) + " does not exist.")
		return
	}

	fmt.Println("room code received: " + string(p) + ", " + string(messageType))

	// Add this new socket to the room-sockets map
	roomConnectionMap[string(p)] = append(roomConnectionMap[string(p)], conn)

	// get list of questions for current room
	QuestionsList := getRoom(string(p))
	// broadcast updated question list to all clients in room
	for _, socket := range roomConnectionMap[string(p)] {
		// send questions DB stuff for code
		err := socket.WriteJSON(QuestionsList)
		if err != nil {
        	fmt.Println("Failed to send through websocket lol. Err: " + err.Error())
			return
		}
		fmt.Println("Broadcasting questions for room " + string(p) + ".")
	}
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	// unmarhall the question id
	decoder := json.NewDecoder(r.Body)
	req := struct {
		RoomCode   string
		QuestionID string
	}{"", ""}

	err := decoder.Decode(&req)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		return
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

	// get list of questions for current room
	QuestionsList := getRoom(req.RoomCode)
	// broadcast updated question list to all clients in room
	for _, socket := range roomConnectionMap[req.RoomCode] {
		// send questions DB stuff for code
		err := socket.WriteJSON(QuestionsList)
		if err != nil {
			failWithStatusCode(err, "Failed to send through websocket.", w, http.StatusInternalServerError)
			return
		}
		fmt.Println("Broadcasting questions for room " + req.RoomCode + ".")
	}
}

func askQuestionHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	req := struct {
		QuestionText string
		RoomCode     string
	}{"", ""}

	// get question text and room from request
	err := decoder.Decode(&req)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		return
	}

	if profane(req.QuestionText) {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
		fmt.Println("bad word detected: " + req.QuestionText)
		return
	}

	// add new question to DB
	queryString := "INSERT INTO questions(room_code, text, votes) VALUES($1, $2, 0)"
	stmt, err := db.Prepare(queryString)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusInternalServerError), w, http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(req.RoomCode, req.QuestionText)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusInternalServerError), w, http.StatusInternalServerError)
		return
	}
	fmt.Printf("code: %s, text: %s\n", req.RoomCode, req.QuestionText)

	// get list of questions for current room
	QuestionsList := getRoom(req.RoomCode)
	// broadcast updated question list to all clients in room
	for _, socket := range roomConnectionMap[req.RoomCode] {
		// send questions DB stuff for code
		err := socket.WriteJSON(QuestionsList)
		if err != nil {
			failWithStatusCode(err, "Failed to send through websocket.", w, http.StatusInternalServerError)
			return
		}
		fmt.Println("Broadcasting questions for room " + req.RoomCode + ".")
	}
}

func hideHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	req := struct {
		RoomCode   string
		QuestionID string
		Secret     string
	}{"", "", ""}

	err := decoder.Decode(&req)

	if err != nil {
		failWithStatusCode(err, http.StatusText(http.StatusBadRequest), w, http.StatusBadRequest)
	}

	queryString := "UPDATE questions SET hide = NOT hide" //toggle hidden status
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

	// get list of questions for current room
	QuestionsList := getRoom(req.RoomCode)
	// broadcast updated question list to all clients in room
	for _, socket := range roomConnectionMap[req.RoomCode] {
		// send questions DB stuff for code
		err := socket.WriteJSON(QuestionsList)
		if err != nil {
			failWithStatusCode(err, "Failed to send through websocket.", w, http.StatusInternalServerError)
			return
		}
		fmt.Println("Broadcasting questions for room " + req.RoomCode + ".")
	}

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
	if isHeroku {
		db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	}
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
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	fmt.Printf("Listening on port: %s\n", port)

	fs := http.FileServer(http.Dir("dist"))
	http.Handle("/", fs)
	http.HandleFunc("/createRoom", createRoomHandler)
	http.HandleFunc("/joinRoom", joinRoomHandler)
	http.HandleFunc("/askQuestion", askQuestionHandler)
	http.HandleFunc("/vote", voteHandler)
	http.HandleFunc("/hide", hideHandler)
	http.ListenAndServe(port, nil)
}
