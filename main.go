package main

import (
	"fmt"
	"log"
	"strings"
	"path"
	"io/ioutil"
	"net"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type StaticHandler struct {
	http.Dir
}

var (
	answer = ""
	webRoot = "web-root"
	whitelistFile = "private/whitelist"
)

func (sh *StaticHandler) ServeHttp(responseWriter http.ResponseWriter, request *http.Request) {
	staticFilePath := staticFilePath(request)
	
	fileHandle, error := sh.Open(staticFilePath)
	if serve404OnError(error, responseWriter) { return }
	defer fileHandle.Close()
	
	fileInfo, error := fileHandle.Stat()
	if serve404OnError(error, responseWriter) { return }
	
	if fileInfo.IsDir() {
		fileHandle, error = sh.Open(staticFilePath + "index.html")
		if serve404OnError(error, responseWriter) { return }
		defer fileHandle.Close()
		
		fileInfo, error = fileHandle.Stat()
		if serve404OnError(error, responseWriter) { return }
	}
	
	http.ServeContent(responseWriter, request, fileInfo.Name(), fileInfo.ModTime(), fileHandle)
}

func loadAnswer() {
	loadedAnswer, error := ioutil.ReadFile("private/answer.txt")
	panicOnError(error)
	answer = string(loadedAnswer)
}

func serveStaticFilesOr404(responseWriter http.ResponseWriter, request *http.Request) {
	staticHandler := StaticHandler{http.Dir(webRoot)}
	staticHandler.ServeHttp(responseWriter, request)
}

func serve404OnError(error error, responseWriter http.ResponseWriter) bool {
	if error != nil {
		responseWriter.WriteHeader(http.StatusNotFound)
		errorTemplate404Content, error := ioutil.ReadFile(webRoot + "/error-templates/404.html")
		panicOnError(error)
		fmt.Fprint(responseWriter, string(errorTemplate404Content))
		return true
	}
	return false
}

func staticFilePath(request *http.Request) string {
	staticFilePath := request.URL.Path
	if !strings.HasPrefix(staticFilePath, "/") {
		staticFilePath = "/" + staticFilePath
		request.URL.Path = staticFilePath
	}
	return path.Clean(staticFilePath)
}

func panicOnError(error error) { if error != nil { panic(error) } }

func answerPost(responseWriter http.ResponseWriter, request *http.Request, requestParameters httprouter.Params) {
	if answer == "" || answer != request.PostFormValue("answer") {
		fmt.Fprint(responseWriter, "false")
		return
	}
	
	ip, _, error := net.SplitHostPort(request.RemoteAddr)
	if error != nil {
		fmt.Fprint(responseWriter, "false")
		return
	}
	fmt.Println(net.ParseIP(ip))
	return
	error = ioutil.WriteFile(whitelistFile, []byte(net.ParseIP(ip)), 0000)
	if error != nil {
		fmt.Fprint(responseWriter, "false")
		return
	}
	fmt.Fprint(responseWriter, "true")
}

func main() {
	loadAnswer()
	router := httprouter.New()
	router.POST("/answer", answerPost)
	router.NotFound = http.HandlerFunc(serveStaticFilesOr404)
	log.Fatal(http.ListenAndServe(":8090", router))
}