//a example to mock how to handle global error
package main

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
)

func main() {
	//Configuration routing
	m := mux.NewRouter()
	m.Handle("/", recoverWrap(http.HandlerFunc(handler))).Methods("GET")
	//It has a panic in the function named handler
	//we will catch the panic when we visit the "http://localhost:8080"

	http.Handle("/", m)
	fmt.Println("Listening...")

	http.ListenAndServe(":8080", nil)

}

//mock a panic happened in the function
func handler(w http.ResponseWriter, r *http.Request) {
	panic(errors.New("panicing from error"))
}

//It is a wrap function to catch panic in the  wrapped  function
func recoverWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					//if it is string,
					//we will wrap it to a error to output
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				fmt.Println("get the panic")
				//if get the panic
				//we will output the error info
				//and return the http errorStatus
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
