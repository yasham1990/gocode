package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/drone/routes"
)

func main() {
	mux := routes.New()
	mux.Post("/mypost", DoPost)
	http.Handle("/", mux)
	fmt.Print("Listening...")
	http.ListenAndServe(":3000", nil)
}

func DoPost(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println("Body: " + string(body[:]))
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(`{"hello":"done"}`))
}
