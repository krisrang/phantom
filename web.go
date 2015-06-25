package main

import (
  "fmt"
  "net/http"
  "os"
)

func main() {
  http.HandleFunc("/", test)
  fmt.Println("listening...")
  err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
  if err != nil {
    panic(err)
  }
}

func test(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintln(res, "hello, heroku")
}
