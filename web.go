package main

import (
  "fmt"
  "net/http"
  "os/exec"
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
  path, _ := exec.LookPath("phantomjs")
  fmt.Fprintln(res, path)
  
  out, _ := exec.Command(path, "-v").Output()
  fmt.Fprintln(res, out)
}
