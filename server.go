package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

//Структура JSON в POST запросе
//{"Values":
//	[
//		{"FirstName": "", "LastName": "", "Age": ""},
//		{"FirstName": "", "LastName": "", "Age": ""}
//	]
//}

var requestData = &requestJson{}

var names struct {
	id int
	fName string
	lName string
	age string
}

type requestJson struct {
	Values []dbStruct `json:"Values"`
}

type dbStruct struct {
	FirstName string `json:"FirstName"`
	LastName string `json:"LastName"`
	Age string `json:"Age"`
}

func parsingPost(r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData)
	_ = r.Body.Close()
	if err != nil {
		fmt.Println("error:", err)
	}
}

func dbConnect() *sql.DB {
	db, err := sql.Open("mysql", "???")
	errPrinter(err)
	return db
}

func errPrinter(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func answer(w http.ResponseWriter, r *http.Request) {

	// Проверка метода запроса
	db := dbConnect()

	if r.Method == http.MethodPost {

		parsingPost(r)

		for _, i := range requestData.Values {
			query1 := "INSERT INTO local_db.first (`id`, `Имя`, `Фамилия`, `Возраст`) " +
				"VALUES (NULL, '" + i.FirstName + "', '" + i.LastName + "', '" + i.Age + "');"
			insert, err := db.Query(query1)
			errPrinter(err)
			_ = insert.Close()
		}

		_, _ = w.Write([]byte("Успешно"))

	} else if r.Method == http.MethodGet {

		query2 := "SELECT * FROM local_db.first"
		sel, err := db.Query(query2)
		errPrinter(err)

		var ans string

		for sel.Next() {
			_ = sel.Scan(&names.id, &names.fName, &names.lName, &names.age)
			ans += strconv.Itoa(names.id) + " ; " + names.fName + " ; " + names.lName + " ; " + names.age + "\n"
		}

		_ = sel.Close()

		if ans != "" {
			resp, _ := json.Marshal(ans)
			_, _ = w.Write(resp)

		} else {
			_, _ = w.Write([]byte("Пустая таблица"))
		}

	} else {
		_, _ = w.Write([]byte(r.Method + " запрос не поддерживается."))
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", answer)
	err := http.ListenAndServe(":3000", mux)
	println(err)
}
