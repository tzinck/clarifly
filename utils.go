package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func failWithStatusCode(err error, msg string, w http.ResponseWriter, statusCode int) {
	failGracefully(err, msg)
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, msg)
}

func failGracefully(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
		panic(err)
	}
}

func randString(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func getRoom(code string) Room {
	queryString := "SELECT room_code, start_time FROM rooms WHERE room_code = $1"
	stmt, err := db.Prepare(queryString)

	failGracefully(err,"Could not prepare query\n")

	var room Room
	err = stmt.QueryRow(code).Scan(&room.Code, &room.Time)

	failGracefully(err,"Could not query\n")

	queryString2 := "SELECT q_id, text, votes, reports, hide, ask_time FROM questions WHERE room_code = $1"
	stmt, err = db.Prepare(queryString2)
	failGracefully(err,"Could not prepare query2\n")
	rows, err := stmt.Query(code)
	failGracefully(err,"Could not query2\n")

	defer rows.Close()

	for rows.Next() {
		var q Question
		rows.Scan(&q.QID, &q.Text, &q.Votes, &q.Reports, &q.Hidden, &q.Time)
		room.Questions = append(room.Questions, q)
	}

	return room
}

func checkHeroku() bool {
	if os.Getenv("IS_HEROKU") != "" {
		fmt.Printf("this is running on heroku")
		return true
	}
	return false
}
