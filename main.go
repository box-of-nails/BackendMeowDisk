package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "URL:", r.URL.String())
}
func loginPage(w http.ResponseWriter, r*http.Request) {

	//log.Println(string(body))
	var users user_data
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &users)
	//expiration:= time.Now().Add(10*time.Hour)
	//cookie:=http.Cookie{
	//	Name: "session_id",
	//	Value: "MMRN9FDZx02MMgVo",
	//	Expires: expiration,
	//	HttpOnly: true,
	//}
	//http.SetCookie(w,&cookie)
	//http.Redirect(w,r, "/",http.StatusFound)
	w.Write([]byte {'h','e','l'})
}
func logoutPage(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}


	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	//http.Redirect(w, r, "/", http.StatusFound)
}
type user_data struct {
	Id    string    `json:"id"`
	Login string `json:"login"`
	Password string `json:"password"`
}

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	pass = "12345"
	dbname = "postgres"
)
func reqistration(w http.ResponseWriter,r*http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var users user_data
	err = json.Unmarshal(body, &users)
	check_put:=put_database(&users)
	if (!check_put){
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)

	if err != nil {
		panic(err)
	}
}

func put_database(data *user_data) bool {

	psqlconn:= fmt.Sprintf("host=%s port=%d user= %s password=%s dbname=%s sslmode=disable",host,port,user,pass,dbname)

	db, err := sql.Open("postgres",psqlconn )

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//var users user_data

	_, err = db.Exec(`
	insert into user_data (
	                     id,
	                     login,
	                     password
	                     )
	values ($1,$2,$3)`,
		data.Id,
		data.Login,
		data.Password,
	)
	if err!=nil{
		return false
	}
	return true

}

//func search_in_database(data *user_data){
//	psqlconn:= fmt.Sprintf("host=%s port=%d user= %s password=%s dbname=%s sslmode=disable",host,port,user,pass,dbname)
//
//	db, err := sql.Open("postgres",psqlconn )
//
//	err = db.Ping()
//	if err != nil {
//		panic(err)
//	}
//	rows, err = db.Query(`SELECT "id","login","password" FROM "user_data"`)
//	if err!=nil{
//		panic(err)
//	}
//	defer rows.Close()
//	for rows.Next(){
//
//		rows.Scan(&id,&login,&password)
//
//	}
//
//
//
//
//}


func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/login",loginPage)
	mux.HandleFunc("/logout",logoutPage) //  потом
	mux.HandleFunc("/register",reqistration) // get json

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("starting server at :8080")
	server.ListenAndServe()
}
