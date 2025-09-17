package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const(
  host = "localhost"
  port = 5432
  user = "postgres"
  password = "postgres"
  dbname = "library"
)

func Connect() *sql.DB{
	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", connInfo)
	if err != nil{
		fmt.Println("Не удалось подключиться к бд")
		panic(err)
	}

	err = db.Ping()
	if err != nil{
		panic(err)
	}
	fmt.Println("Подключение к бд успешно")
	return db
}

func CloseConnection(db *sql.DB){
	db.Close()
}
