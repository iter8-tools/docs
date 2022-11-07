package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	MY_VERSION         = "MY_VERSION"
	DEFAULT_VERSION = "v1"
)

func getVersion() string {
	version, ok := os.LookupEnv(MY_VERSION)
	if !ok {
		version = DEFAULT_VERSION
	}
	return version
}

// implment /recommend endpoint returning value of VERSION env variable
func recommend(w http.ResponseWriter, req *http.Request) {
	Logger.Trace("recommend called")
	version := getVersion()
	Logger.Info("/recommend returns ", version)
	fmt.Fprintln(w, version)
}

var Logger *logrus.Logger

func main() {
	Logger = logrus.New()

	// configure backend service with "/recommend" endpoint
	http.HandleFunc("/recommend", recommend)
	http.ListenAndServe(":8091", nil)
}
