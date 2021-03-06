package main

import (
	"api-test/database"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type userData struct {
	Login    string
	Password string
}

type taskData struct {
	Id		int
	Title   string
	Goal    string
	Duedate string
}

func handleTaskData(r *http.Request) *taskData {
	var data taskData

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}
	return &data
}

func handleUserData(r *http.Request) *userData {
	var data userData

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}
	return &data
}

func getPasswordByLogin(login string) string {
	sqlString := "SELECT password FROM users WHERE login=$1;"
	row := database.GetSpecificFromDB(sqlString, login)
	return row
}

func hashPassword(password string) string{
	hash := sha1.New()
	hash.Write([]byte(password))
	sum := hex.EncodeToString(hash.Sum(nil))
	return sum
}

func addTasks(w http.ResponseWriter, r *http.Request) {
	data := handleTaskData(r)
	sqlString := "INSERT INTO tasks(title, goal, duedate) VALUES($1, $2, $3) returning id;"
	values := []string {data.Title, data.Goal, data.Duedate}
	database.InsertToDB(sqlString, values)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "ok"}`))
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	sqlString := "SELECT * FROM tasks;"
	tasks := database.GetAllFromDB(sqlString)
	JSON, _ := json.Marshal(tasks)
	w.WriteHeader(http.StatusCreated)
	w.Write(JSON)
}

func deleteTasks(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	sqlString := "DELETE FROM tasks WHERE id = $1;"
	database.DeleteFromDB(sqlString, id)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "ok"}`))
}

func updateTasks(w http.ResponseWriter, r *http.Request) {
	data := handleTaskData(r)
	sqlString := "UPDATE tasks SET title = $1, goal = $2, duedate = $3 WHERE id = $4;"
	id := strconv.Itoa(data.Id)
	values := []string {data.Title, data.Goal, data.Duedate, id}
	database.UpdateInDB(sqlString, values)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "ok"}`))
}

func signIn(w http.ResponseWriter, r *http.Request) {
	data := handleUserData(r)
	passwordForChecking := getPasswordByLogin(data.Login)
	if passwordForChecking == "" {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message": "error"}`))
		return
	}
	hashedPassword := hashPassword(data.Password)
	if hashedPassword != passwordForChecking {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message": "error"}`))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "ok"}`))
}

func signUp(w http.ResponseWriter, r *http.Request) {
	data := handleUserData(r)
	isPasswordExists := getPasswordByLogin(data.Login)
	if isPasswordExists != "" {
		fmt.Println("in")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message": "error"}`))
		return
	}
	hashedPassword := hashPassword(data.Password)
	sqlString := "INSERT INTO users(login, password) VALUES($1, $2) returning id;"
	values := []string {data.Login, hashedPassword}
	database.InsertToDB(sqlString, values)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "ok"}`))
}

