package main

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"net/url"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
//	"strings"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

var (
	db *sql.DB
	profileTable = `CREATE TABLE IF NOT EXISTS profile (
		name VARCHAR(64) NULL DEFAULT NULL,
		phone VARCHAR(64) NULL DEFAULT NULL,
		password VARCHAR(64) NULL DEFAULT NULL
    );`
)

type Profile struct {
	Name string
	Phone string
	Password struct {
	}
}

func PanicIf(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func setupDB() *sql.DB{
	db, err := sql.Open("mysql", "root@/recovr?charset=utf8")
	PanicIf(err)
	return db
}

func (p Profile) Get(values url.Values) (int, interface {}) {
	fmt.Printf("In Profile GET!")
	data := map[string]string{"hello": "world"}
    return 200, data
}

func (p Profile) Post(values url.Values) (int, interface {}) {
	fmt.Printf("In Profile POST!")
	stmt, err := db.Prepare("INSERT profile SET name=?,phone=?,password=?")
	PanicIf(err)

	name :=  values["name"][0]
	phone := values["phone"][0]
	password := values["password"][0]

	res, err := stmt.Exec(name, phone, password)
	PanicIf(err)
    fmt.Println("Response:", res)

	data := map[string]string{"response": "true", "message": "Profile created successfully" }
	fmt.Println("Data:", data)
	if err != nil {
		return 405, map[string]string{"response": "false", "message": "Something went woring" }
	}

	return 200, data
}

func Abort(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}


func profileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside handler")

	var data interface{}
	var code int
	profile := Profile{}

	r.ParseForm()
	method := r.Method
	values := r.Form

	switch method {
	case GET:
		code, data = profile.Get(values)
	case POST:
		code, data = profile.Post(r.Form)
	default:
		Abort(w, 405)
		return
	}

	content, err := json.Marshal(data)
	if err != nil {
		PanicIf(err)
	}

	w.WriteHeader(code)
	w.Write(content)
}

func main() {
	db = setupDB()
	defer db.Close()

	ctble, err := db.Query(profileTable)
	PanicIf(err)
	fmt.Println("Profile table created successfully:", ctble)


	http.HandleFunc("/profile", profileHandler)

	fmt.Println("Listening on 3000....")
	if err:= http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		log.Fatalf("Error to listen:", err)
	}
}

