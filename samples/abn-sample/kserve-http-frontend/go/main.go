package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	abn "github.com/iter8-tools/iter8/abn/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Logger = logrus.New()

// map of version number to route to backend service
// here route is the URL of the backend service to which the request should be sent
var versionNumberToRoute = []string{
	"backend-0",
	"backend-1",
}

// implment /getRecommendation endpoint
// calls backend service (REST API /recommend endpoint)
func getRecommendation(w http.ResponseWriter, req *http.Request) {
	Logger.Info("/getRecommendation")
	defer Logger.Info("returned ")

	// Get user (session) identifier, for example by inspection of header X-User
	user := req.Header["X-User"][0]

	// Get endpoint of backend endpoint "/recommend"
	// In this example, the backend endpoint depends on the version of the backend service
	// the user is assigned by the Iter8 SDK Lookup() method

	// start with default route
	route := versionNumberToRoute[0]

	// call A/B/n service API Lookup() to get a recommended version for the user
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s, err := getABnClient().Lookup(
		ctx,
		&abn.Application{
			Name: backendName(),
			User: user,
		},
	)
	// if successful, use returned version; otherwise will use the default
	if err != nil {
		Logger.Info("error: " + err.Error())
	}
	// if successful, use recommended version; otherwise will use default
	if err == nil && s != nil {
		Logger.Infof("successful call to lookup %d", s.GetVersionNumber())
		versionNumber := int(s.GetVersionNumber())
		if err == nil && 0 <= versionNumber && versionNumber < len(versionNumberToRoute) {
			route = versionNumberToRoute[versionNumber]
		} // else use default value for route
	}
	Logger.Info("lookup suggested route " + route)

	// call backend
	resp, err := callBackend(route)
	if err != nil {
		Logger.Errorf("call to backend failed: %s", err.Error())
		http.Error(w, "call to backend failed", http.StatusInternalServerError)
		return
	}

	// write response to query
	fmt.Fprintln(w, "Recommendation: "+resp)
}

// implment /buy endpoint
// writes value for sample_metric which may have spanned several calls to /getRecommendation
// in this sample, the metric value is random
func buy(w http.ResponseWriter, req *http.Request) {
	Logger.Info("/buy")
	defer Logger.Info("returned ", http.StatusAccepted)
	// Get user (session) identifier, for example by inspection of header X-User
	user := req.Header["X-User"][0]

	// export metric to metrics database; this is best effort; ignore any failure
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, _ = getABnClient().WriteMetric(
		ctx,
		&abn.MetricValue{
			Name:        "sample_metric",
			Value:       fmt.Sprintf("%f", rand.Float64()*100.0), // strconv.Itoa(rand.Intn(100)),
			Application: backendName(),
			User:        user,
		},
	)
	fmt.Fprintln(w, "Purchase complete")
}

func main() {
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		DisableQuote:    true,
		DisableSorting:  true,
	})
	Logger.SetLevel(logrus.TraceLevel)

	// configure frontend service with "/hello" and "/goodbye" endpoints
	http.HandleFunc("/getRecommendation", getRecommendation)
	http.HandleFunc("/buy", buy)
	http.ListenAndServe(":8090", nil)
}

var abnClient *abn.ABNClient

func getABnClient() abn.ABNClient {
	if abnClient == nil {
		// establish connection to A/B/n service
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		conn, err := grpc.Dial(
			fmt.Sprintf(
				"%s:%s",
				lookupEnv("ABN_SERVICE", "iter8"),
				lookupEnv("ABN_SERVICE_PORT", "50051"),
			),
			opts...,
		)
		if err != nil {
			panic("Cannot establish connection with abn service")
			// return
		}
		c := abn.NewABNClient(conn)
		abnClient = &c

	}

	return *abnClient
}

func lookupEnv(key string, default_value string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return default_value
}

func backendName() string {
	return lookupEnv("BACKEND_APPLICATION_NAME", "default/backend")
}

// callBackend calls infer endpoint using HTTP
// equivalent to:
//
//	curl -X POST http://backend-0.default.svc.cluster.local/v2/models/sklearn-irisv2/infer' \
//	   -H 'Content-Type: application/json' \
//	   --data data
func callBackend(route string) (string, error) {
	data := strings.Replace(`{
		"inputs": [
		  {
			"name": "input-0",
			"shape": [2, 4],
			"datatype": "FP32",
			"data": [
			  [6.8, 2.8, 4.8, 1.4],
			  [6.0, 3.4, 4.5, 1.6]
			]
		  }
		]
	  }`, "\n", "", -1)
	resp, err := http.Post(
		fmt.Sprintf("http://%s.default.svc.cluster.local/v2/models/%s/infer", route, route),
		"application/json",
		bytes.NewBuffer([]byte(data)),
	)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
