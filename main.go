package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"database/sql"
	"fmt"
	"encoding/json"
	"strconv"
	"io/ioutil"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/godror/godror"
)

type Note struct{
	Id int				`json:"id"`
	Title string		`json:"title"`	
	Description string	`json:"description"`
}


func GetConnection() *sql.DB {
	var db *sql.DB
	db, errdb := sql.Open("godror", "NOTESUSER/admin@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=tcp)(HOST=localhost)(PORT=1521)))(CONNECT_DATA=(SID=xe)))")
	if errdb != nil {
		fmt.Println(errdb)
		return nil
	}
	//defer db.Close()
	return db
}

func holaName(w http.ResponseWriter, r *http.Request){
	name := mux.Vars(r)["name"]
	fmt.Fprint(w, "hola "+name)
}

func getNotes(w http.ResponseWriter, r *http.Request){
	db := GetConnection()
	var notes []Note
	rows, _ := db.Query("SELECT id, title, description FROM notes")
    var note Note
    for rows.Next() {
		rows.Scan(&note.Id, &note.Title, &note.Description)
		notes=append(notes,note)
    }
	json.NewEncoder(w).Encode(notes)
}

func getNote(w http.ResponseWriter, r *http.Request){
	db := GetConnection()
	defer db.Close()
	id_String := mux.Vars(r)["id"]
	var note Note
	var error1 string
	var error2 string
	db.Exec("BEGIN getNoted (:1,:2,:3,:4,:5); end;",id_String, sql.Out{Dest:&note.Title}, sql.Out{Dest:&note.Description}, sql.Out{Dest:&error1}, sql.Out{Dest:&error2})
	if error1 != "" {
		errorMessage := map[string]string{
			"code":"error",
			"message":error2,
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	note.Id,_=strconv.Atoi(id_String)
	json.NewEncoder(w).Encode(note)
}

func createNote(w http.ResponseWriter, r *http.Request){
	db := GetConnection()
	defer db.Close()
	var temp Note
	var error1 string
	var error2 string
	body,_:=ioutil.ReadAll(r.Body)
	json.Unmarshal(body,&temp)
	db.Exec("BEGIN NewNote (:1,:2,:3,:4); end;",temp.Title, temp.Description,sql.Out{Dest:&error1},sql.Out{Dest:&error2})
	if error1 != "" {
		errorMessage := map[string]string{
			"code":"error",
			"message":error2,
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	
	message1 := map[string]string{
		"code":"correcto",
		"message":"Nota registrada correctamente",
	}
	json.NewEncoder(w).Encode(message1)
}

func updateNote(w http.ResponseWriter, r *http.Request){
	db := GetConnection()
	var temp Note
	var error1 string
	var error2 string
	id_String := mux.Vars(r)["id"]
	body,_:=ioutil.ReadAll(r.Body)
	json.Unmarshal(body,&temp)
	db.Exec("BEGIN UpdateNote (:1,:2,:3,:4,:5); end;",id_String,temp.Title, temp.Description,sql.Out{Dest:&error1},sql.Out{Dest:&error2})
	if error1 != "" {
		errorMessage := map[string]string{
			"code":"error",
			"message":error2,
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	
	message1 := map[string]string{
		"code":"correcto",
		"message":"Nota actualizada correctamente",
	}
	json.NewEncoder(w).Encode(message1)
}

func deleteNote(w http.ResponseWriter, r *http.Request){
	db := GetConnection()

	var error1 string
	var error2 string
	id_String := mux.Vars(r)["id"]
	db.Exec("BEGIN DeleteNote (:1,:2,:3); end;",id_String,sql.Out{Dest:&error1},sql.Out{Dest:&error2})
	if error1 != "" {
		errorMessage := map[string]string{
			"code":"error",
			"message":error2,
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	
	message1 := map[string]string{
		"code":"correcto",
		"message":"Nota actualizada correctamente",
	}
	json.NewEncoder(w).Encode(message1)
}

func main(){

	r := mux.NewRouter()
	r.HandleFunc("/note", createNote).Methods("POST")
	r.HandleFunc("/notes", getNotes).Methods("GET")
	r.HandleFunc("/notes/{id}", getNote).Methods("GET")
	r.HandleFunc("/notes/{id}", updateNote).Methods("PUT")
	r.HandleFunc("/notes/{id}", deleteNote).Methods("DELETE")

	log.Print("Corriendo en el puerto 8085")
	err := http.ListenAndServe(":8085",r)
	if err != nil{
		log.Fatal("error: ",err)
	}
}