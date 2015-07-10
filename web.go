package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

var scriptPath string
var port int

const (
	VERSION      = "0.2.0"
	VERSIONFANCY = "Choosy Phantom"
)

func init() {
	flag.StringVar(&scriptPath, "script", "", "path to phantomjs/slimerjs script file")
	flag.IntVar(&port, "port", 5000, "port to run phantom on")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	checkExecutable("phantomjs")
	checkExecutable("slimerjs")
	checkFile(scriptPath)

	portString := strconv.Itoa(port)

	http.HandleFunc("/screenshot", screenshot)
	http.HandleFunc("/version", version)
	fmt.Printf("Phantom listening on port %s\n", portString)
	err := http.ListenAndServe(":"+portString, nil)
	if err != nil {
		panic(err)
	}
}

func version(res http.ResponseWriter, req *http.Request) {
	phantomOut, _ := exec.Command("phantomjs", "-v").Output()
	slimerOut, _ := exec.Command("slimerjs", "-v").Output()
	fmt.Fprintf(res, "Phantom version %s %s\n", VERSION, VERSIONFANCY)
	fmt.Fprintf(res, "Phantom version is %s", phantomOut)
	fmt.Fprintf(res, "SlimerJS version is %s\n", slimerOut)
}

func screenshot(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	url := req.Form.Get("url")
	if url == "" {
		fmt.Fprintln(res, "Error: must specify URL")
		return
	}

	engine := req.Form.Get("engine")
	if engine == "" || (engine != "phantomjs" && engine != "slimerjs") {
		engine = "slimerjs"
	}

	file, _ := ioutil.TempFile(os.Getenv("TMPDIR"), "screenshot")
	tmpPath := file.Name() + ".png"
	file.Close()
	os.Remove(file.Name())
	defer os.Remove(tmpPath)

	out, cmdErr := cmd(engine, scriptPath, url, tmpPath).CombinedOutput()
	if cmdErr != nil {
		fmt.Fprintf(res, "Error executing capture engine: %s\n", out)
		return
	}

	http.ServeFile(res, req, tmpPath)
}

func cmd(executable string, args ...string) *exec.Cmd {
	path, _ := exec.LookPath(executable)
	out := exec.Command(path, args...)
	return out
}

func checkExecutable(name string) {
	_, err := exec.LookPath(name)
	if err != nil {
		log.Fatalf("Installing %s is in your future", name)
	}
}

func checkFile(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("File %s doesn't exist", name)
		}
	}
	return true
}
