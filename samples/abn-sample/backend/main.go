package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type Data struct {
	Id     int
	Name   string
	Source string
}

// implment /recommend endpoint returning value of VERSION env variable
func recommend(w http.ResponseWriter, req *http.Request) {
	Logger.Trace("recommend called")

	data := Data{
		Id:     19,
		Name:   "sample",
		Source: os.Getenv("HOSTNAME"),
	}
	json.NewEncoder(w).Encode(data)
}

var Logger *logrus.Logger

func main() {
	Logger = logrus.New()

	// configure backend service with "/recommend" endpoint
	http.HandleFunc("/recommend", recommend)
	http.ListenAndServe(":8091", nil)
}
