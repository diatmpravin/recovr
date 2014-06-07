package main

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"net/url"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
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

func (p Profile) Get(values url.Values) (int, interface {}) {
	data := map[string]string{"hello": "world"}
    return 200, data
}

func Abort(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}


func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside handler")

	var data interface{}
	var code int
	profile := Profile{}

	method := r.Method
	values := r.Form
	fmt.Println("Method: ", method)
	fmt.Println("Values: ", values)

	switch method {
	case GET:
		code, data = profile.Get(values)
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
	http.HandleFunc("/register", registerHandler)

	fmt.Println("Listening on 3000....")
	if err:= http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		log.Fatalf("Error to listen:", err)
	}
}

