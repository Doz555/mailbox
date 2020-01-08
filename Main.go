package main

//import . "llrp"
//import "os"
//import "io"
import "fmt"

import _ "github.com/go-sql-driver/mysql"
import "github.com/gorilla/mux"
import "log"
import "net/http"

import "time"





func main() {
	start := time.Now()
	defer func() {
		fmt.Printf("%s", time.Now().Sub(start).String)
	}()

	r := mux.NewRouter()

	r.PathPrefix("/src/").Handler(http.StripPrefix("/src/", http.FileServer(http.Dir("./src/"))))
	r.PathPrefix("/pic/").Handler(http.StripPrefix("/pic/", http.FileServer(http.Dir("./pic/"))))
	r.PathPrefix("/bower_components/").Handler(http.StripPrefix("/bower_components/", http.FileServer(http.Dir("./bower_components/"))))

	// Bind to a port and pass our router in
	//log.Fatal(http.ListenAndServe(":8000", r))
	fmt.Printf("start server 80 ")
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
