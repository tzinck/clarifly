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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Connection success!")
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

    // message := struct {
    //  RoomString string
    // }{""}
    // frontend handshake to get user and hook them into the userMap for sockets
    messageType, p, err := conn.ReadMessage()
    if err != nil {
        failWithStatusCode(err, "Failed to handshake", w, http.StatusInternalServerError)
        return
    }

    fmt.Println("room code received: " + string(p))
    fmt.Println(messageType)
    roomConnectionMap[string(p)] = append(roomConnectionMap[string(p)], conn)

    returnmsg := []byte("we got yo shit")
    err = conn.WriteMessage(1, returnmsg)
    if err != nil {
        fmt.Println("bad lol")
        return
    }
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
		fmt.Println("11111")
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

	// update them sockets
	room := getRoom(req.RoomID)
	for _, socket := range roomConnectionMap[req.RoomID] {
		// grab all the questions from the database for this room and send them back over the socket
		err = socket.WriteJSON(room)
		if err != nil {
			failWithStatusCode(err, http.StatusText(http.StatusInternalServerError), w, http.StatusInternalServerError)
			return
		}
	}
}

func askQuestionHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
    req := struct {
        QuestionText string
        RoomCode string
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
        err = socket.WriteJSON(QuestionsList)
        if err != nil {
	        failWithStatusCode(err, "Failed to send through websocket.", w, http.StatusInternalServerError)
	        return
	    }
        //fmt.Println(ws)
        fmt.Println("sent the questions to a guy\n");
    }
}

func hideHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	req := struct {
		RoomID     string
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

	room := getRoom(req.RoomID)
	for _, socket := range roomConnectionMap[req.RoomID] {
		// grab all the questions from the database for this room and send them back over the socket
		err = socket.WriteJSON(room)
		if err != nil {
			failWithStatusCode(err, http.StatusText(http.StatusInternalServerError), w, http.StatusInternalServerError)
			return
		}
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
	http.HandleFunc("/askQuestion", askQuestionHandler)
	http.HandleFunc("/vote", voteHandler)
	http.HandleFunc("/hide", hideHandler)
	http.ListenAndServe(configuration.Port, nil)
}
