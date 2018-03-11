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
