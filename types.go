package main

type QuestionMessage struct {
	Question string
	Votes    int
	Q_ID     int
}

type Configuration struct {
	DB DbCreds
}

type DbCreds struct {
	Host   string
	Port   int
	User   string
	Pass   string
	DbName string
}

type Question struct {
	QID     int
	Text    string
	Votes   int
	Reports int
	Hidden  bool
	Time    string
}

type Room struct {
	Code      string
	Time      string
	Questions []Question
}
