package main

import (
    "fmt"
    "log"
    "net/http"
    "transfer-service/handlers"

    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/transfers", handlers.CreateTransfer).Methods("POST")
    r.HandleFunc("/transfers", handlers.GetTransfers).Methods("GET")

    fmt.Println("Transfer service running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
