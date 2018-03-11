package main


type QuestionMessage struct {
	Question string
    Votes int
    Q_ID int
}

type Configuration struct {
	DB   DbCreds
	Port string
}

type DbCreds struct {
	Host   string
	Port   int
	User   string
	Pass   string
	DbName string
}

type Question struct {
	QID         int
	UID         int
	Text        string
	Votes       int
	Reports     int
	Hidden      bool
	Time        string
}

type Room struct {
	Code        string
	Creator     int
	Time        string
	Questions []Question
}
