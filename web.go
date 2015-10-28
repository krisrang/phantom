package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

var scriptPath string
var port int

const (
	VERSION      = "0.3.0"
	VERSIONFANCY = "Diffy Phantom"
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

	fmt.Printf("Snapping %v\n", url)

	engine := req.Form.Get("engine")
	if engine == "" || (engine != "phantomjs" && engine != "slimerjs") {
		engine = "slimerjs"
	}

	zoom := req.Form.Get("zoom")
	if zoom == "" {
		zoom = "1"
	}

	diff := req.Form.Get("diff")

	file, _ := ioutil.TempFile(os.Getenv("TMPDIR"), "screenshot")
	tmpPath := file.Name() + ".png"
	file.Close()
	os.Remove(file.Name())
	defer os.Remove(tmpPath)

	if diff == "true" {
		tmp1Path := file.Name() + "1.png"
		tmp2Path := file.Name() + "2.png"
		defer os.Remove(tmp1Path)
		defer os.Remove(tmp2Path)

		var wg sync.WaitGroup
		wg.Add(2)
		go snapUrl(&wg, engine, scriptPath, url+"&pass=1", zoom, tmp1Path)
		go snapUrl(&wg, engine, scriptPath, url+"&pass=2", zoom, tmp2Path)
		wg.Wait()

		out, cmdErr := cmd("convert", tmp1Path, tmp2Path, "-alpha", "off", "(", "-clone", "0,1", "-compose", "difference", "-composite", "-separate", "-evaluate-sequence", "max", "-auto-level", "-negate", ")", "(", "-clone", "0,2", "-fx", "v==0?0:u/v-u.p{0,0}/v+u.p{0,0}", ")", "-delete", "0,1", "+swap", "-compose", "Copy_Opacity", "-composite", tmpPath).CombinedOutput()
		if cmdErr != nil {
			fmt.Fprintf(res, "Error processing image: %s\n", out)
			return
		}
	} else {
		url = base64.StdEncoding.EncodeToString([]byte(url))
		out, cmdErr := cmd(engine, "--disk-cache=true", "--max-disk-cache-size=102400", scriptPath, url, zoom, tmpPath).CombinedOutput()
		if cmdErr != nil {
			fmt.Fprintf(res, "Error executing capture engine: %s\n", out)
			return
		}

		fmt.Println(string(out))
	}

	http.ServeFile(res, req, tmpPath)
}

func cmd(executable string, args ...string) *exec.Cmd {
	path, _ := exec.LookPath(executable)
	out := exec.Command(path, args...)
	return out
}

func snapUrl(wg *sync.WaitGroup, engine string, scriptPath string, url string, zoom string, tmpPath string) ([]byte, error) {
	defer wg.Done()

	out, cmdErr := cmd(engine, scriptPath, url, zoom, tmpPath).CombinedOutput()
	if cmdErr != nil {
		return nil, fmt.Errorf("Error executing capture engine: %s\n", out)
	}

	return out, cmdErr
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
