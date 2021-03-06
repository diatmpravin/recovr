package main

// import required libraries
import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"net/url"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"time"
)
// global constants
const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

// global variable
var (
	db *sql.DB
	profileTable = `CREATE TABLE IF NOT EXISTS profile (
		name VARCHAR(64) NULL DEFAULT NULL,
		phone VARCHAR(64) NULL DEFAULT NULL,
		password VARCHAR(64) NULL DEFAULT NULL
    );`

	coordinateTable = `CREATE TABLE IF NOT EXISTS coordinate (
		phone VARCHAR(64) NULL DEFAULT NULL,
		lat VARCHAR(64) NULL DEFAULT NULL,
		lng VARCHAR(64) NULL DEFAULT NULL,
		createAt VARCHAR(64) NULL DEFAULT NULL
    );`
)

type Date int64

type Login struct {
	Phone string
	Password string
}

type Profile struct {
	Login    //anonymous field
	Password string
}

type Coordinate struct {
	Phone string
	Lat string
	Lng string
	CreateAt Date
}

// run before program execution
func init() {
	db = setupDB()

	ptable, err := db.Query(profileTable)
	PanicIf(err)
	fmt.Println("Profile table created successfully:", ptable)
	ctable, err := db.Query(coordinateTable)
	PanicIf(err)
	fmt.Println("Coordinate table created successfully:", ctable)
}

// setup db connection
func setupDB() *sql.DB{
	db, err := sql.Open("mysql", "root@/recovr?charset=utf8")
	PanicIf(err)
	return db
}

// global through ERROR
func PanicIf(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

// Abort with ERROR code
func Abort(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
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

func (p Profile) Put(values url.Values) (int, interface {}) {
	fmt.Printf("In Profile PUT!")

	password := values["password"][0]

	var s string
	err := db.QueryRow("SELECT phone FROM profile WHERE password=?", password).Scan(&s)
	PanicIf(err)
	fmt.Println("S", s)

	data := map[string]string{"response": "true", "message": "Logged in successfully" }
	fmt.Println("Data:", data)
	if err != nil {
		return 405, map[string]string{"response": "false", "message": "Either your user or password is incorrect!" }
	}

	return 200, data
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside handler")

	var data interface{}
	var code int
	profile := Profile{}

	r.ParseForm()
	method := r.Method
	fmt.Println("Method:", method)
	values := r.Form

	switch method {
	case GET:
		code, data = profile.Get(values)
	case POST:
		code, data = profile.Post(r.Form)
	case PUT:
		code, data = profile.Put(r.Form)
	default:
		Abort(w, 405)
		return
	}

	content, err := json.Marshal(data)
	if err != nil {
		PanicIf(err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(content)
}

func coordinateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside coordinate Handler")

	var data interface{}
	var code int
	coordinate := Coordinate{}

	r.ParseForm()
	method := r.Method
	values := r.Form

	switch method {
	case GET:
		code, data = coordinate.Get(values)
	case POST:
		code, data = coordinate.Post(values)
	default:
		Abort(w, 405)
		return
	}

	content, err := json.Marshal(data)
	if err != nil {
		PanicIf(err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(content)
}

func (c Coordinate) Post(values url.Values) (int, interface {}) {
	fmt.Printf("In Coordinate POST!")
	stmt, err := db.Prepare("INSERT coordinate SET phone=?,lat=?,lng=?,createAt=?")
	PanicIf(err)

	phone := values["phone"][0]
	lat :=  values["lat"][0]
	lng :=  values["long"][0]
	createAt := time.Now().Local()

	fmt.Println("Phone: ", phone)
	fmt.Println("Lat: ", lat)
	fmt.Println("lng: ", lng)
	fmt.Println("createAt: ", createAt)

	res, err := stmt.Exec(phone, lat, lng, createAt)
	PanicIf(err)
	fmt.Println("Response:", res)

	data := map[string]string{"response": "true", "message": "Coordinate saved successfully" }
	fmt.Println("Data:", data)
	if err != nil {
		return 405, map[string]string{"response": "false", "message": "Something went woring" }
	}

	return 200, data
}

func (c Coordinate) Get(values url.Values) (int, interface {}) {
	fmt.Printf("In Coordinate GET!")

	limit:=  values["range"][0]
	phone := values["phone"][0]
	fmt.Println("Range: ", limit)
	fmt.Println("Phone: ", phone)

	rows, err := db.Query(" select lat, lng from coordinate where phone=? order by createAt limit ?", phone, limit)
	PanicIf(err)

	coordinates := []Coordinate{}
	for rows.Next() {
		coordinate := Coordinate{}
		err := rows.Scan(&coordinate.Lat, &coordinate.Lng)
		PanicIf(err)
		coordinates = append(coordinates, coordinate)
	}

	fmt.Println("Coordinates ", coordinates)

	data := coordinates
	fmt.Println("Data:", data)
	if err != nil {
		return 405, map[string]string{"response": "false", "message": "Something went woring" }
	}

	return 200, data
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside coordinate Handler")

	var data interface{}
	var code int
	login := Login{}

	r.ParseForm()
	method := r.Method
	values := r.Form

	switch method {
	case POST:
		code, data = login.Post(values)
	default:
		Abort(w, 405)
		return
	}

	content, err := json.Marshal(data)
	if err != nil {
		PanicIf(err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(content)
}

func (p Login) Post(values url.Values) (int, interface {}) {
	fmt.Printf("In Profile PUT!")

	username := values["phone"][0]
	password := values["password"][0]
	fmt.Println("UserName: ", username)
	fmt.Println("Password: ", password)

	var name string
	err := db.QueryRow("SELECT name FROM profile WHERE password=? and phone=?", password, username).Scan(&name)
	PanicIf(err)
	fmt.Println("Name", name)

	data := map[string]string{"response": "true", "message": "Logged in successfully" }
	fmt.Println("Data:", data)
	if err != nil {
		return 405, map[string]string{"response": "false", "message": "Either your user or password is incorrect!" }
	}

	return 200, data
}

func main() {
	defer db.Close()

	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/coordinate", coordinateHandler)

	fmt.Println("Listening on 3000....")
	if err:= http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		log.Fatalf("Error to listen:", err)
	}
}
