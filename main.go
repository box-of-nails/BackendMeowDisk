package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "URL:", r.URL.String())
}
func loginPage(w http.ResponseWriter, r*http.Request) {
	expiration:= time.Now().Add(10*time.Hour)
	cookie:=http.Cookie{
		Name: "session_id",
		Value: "MMRN9FDZx02MMgVo",
		Expires: expiration,
		HttpOnly: true,
	}
	http.SetCookie(w,&cookie)
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
	id int;
	login string
	password int
}

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "12345"
	dbname = "user_data"
)
func reqistration(w http.ResponseWriter,r*http.Request){
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var t user_data
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	log.Println(t)
}


func main() {
	psqlconn:= fmt.Sprintf("host=%s port=%d user= %s password=%s dbname=%s sslmode=disable",host,port,user,password,dbname)

	db, err := sql.Open("postgres",psqlconn )
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`
	insert into user_data (
	                      id,
	                      login,
	                      password
	                      )
	values ($1,$2,$3)`,
	"2",
	"test_log",
	"test_pass",)
	if err!=nil{
		log.Panic(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/login",loginPage)
	mux.HandleFunc("/logout",logoutPage)
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
