package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var (
	answer        = ""
	webRoot       = "web-root"
	whitelistFile = "tmp/whitelist"
)

func main() {
	loadAnswer()
	router := httprouter.New()
	router.POST("/answer", answerPost)
	router.NotFound = http.HandlerFunc(serveStaticFilesOr404)
	log.Fatal(http.ListenAndServe(":8090", router))
}

func answerPost(responseWriter http.ResponseWriter, request *http.Request, requestParameters httprouter.Params) {
	if answer == "" || answer != request.PostFormValue("answer") {
		fmt.Fprint(responseWriter, "false")
		return
	}
	ipAddresses := strings.Split(request.Header.Get("x-forwarded-for"), ", ")
	error := ioutil.WriteFile(whitelistFile, []byte(ipAddresses[0]), 0000)
	if error != nil {
		fmt.Fprint(responseWriter, "false")
		return
	}
	fmt.Fprint(responseWriter, "true")
}
