package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var phantomPath string
var scriptPath string

const (
	VERSION      = "0.1.0"
	VERSIONFANCY = "Snappy Phantom"
)

func main() {
	path, pathErr := exec.LookPath("phantomjs")
	if pathErr != nil {
		log.Fatal("installing phantomjs is in your future")
	}
	phantomPath = path
	scriptPath, _ = filepath.Abs(os.Args[1])

	http.HandleFunc("/screenshot", screenshot)
	http.HandleFunc("/version", version)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func version(res http.ResponseWriter, req *http.Request) {
	out, _ := exec.Command(phantomPath, "-v").Output()
	fmt.Fprintf(res, "Version %s %s\n", VERSION, VERSIONFANCY)
	fmt.Fprintf(res, "Phantom version is %s\n", out)
}

func screenshot(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	url := req.Form.Get("url")
	if url == "" {
		fmt.Fprintln(res, "must specify URL")
		return
	}

	file, _ := ioutil.TempFile(os.Getenv("TMPDIR"), "screenshot")
	tmpPath := file.Name() + ".png"
	file.Close()
	os.Remove(file.Name())

	out, cmdErr := exec.Command(phantomPath, scriptPath, url, tmpPath).CombinedOutput()
	if cmdErr != nil {
		fmt.Fprintf(res, "Error screenshotting page: %s\n", out)
		return
	}
	defer os.Remove(tmpPath)

	http.ServeFile(res, req, tmpPath)
}
