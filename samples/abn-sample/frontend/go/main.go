package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	pb "github.com/iter8-tools/iter8/abn/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// var log *logrus.Logger

var (
	// map of track to route to backend service
	trackToRoute = []string{
		"http://backend.default.svc.cluster.local:8091",
		"http://backend-candidate-1.default.svc.cluster.local:8091",
	}

	// gRPC client connection
	client *pb.ABNClient // set in main()
	Logger = logrus.New()
)

// implment /getRecommendation endpoint
// calls backend service /recommend endpoint
func getRecommendation(w http.ResponseWriter, req *http.Request) {
	Logger.Info("/getRecommendation")
	defer Logger.Info("returned ")

	// Get user (session) identifier, for example by inspection of header X-User
	user := req.Header["X-User"][0]

	// Get endpoint of backend endpoint "/recommend"
	// In this example, the backend endpoint depends on the version (track) of the backend service
	// the user is assigned by the Iter8 SDK Lookup() method

	// start with default route
	route := trackToRoute[0]

	// call ABn service API Lookup() to get an assigned track for the user
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s, err := (*client).Lookup(
		ctx,
		&pb.Application{
			Name: "default/backend",
			User: user,
		},
	)
	if err != nil {
		Logger.Info("error: " + err.Error())
	}
	// if successful, use recommended track; otherwise will use default route
	if err == nil && s != nil {
		Logger.Info("successful call to lookup " + s.GetTrack())
		track, err := strconv.Atoi(s.GetTrack())
		if err == nil && 0 <= track && track < len(trackToRoute) {
			route = trackToRoute[track]
		} // else use default value for route
	}
	Logger.Info("lookup suggested track " + route)

	// call backend service using url
	resp, err := http.Get(route + "/recommend")
	if err != nil {
		http.Error(w, "call to backend endpoint /recommend failed", http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("backend endpoint /recommend returned no data %s", err), http.StatusInternalServerError)
		return
	}

	// write response to query
	fmt.Fprintln(w, "Recommendation: "+string(body))
}

// implment /buy endpoint
// writes value for sample_metric which may have spanned several calls to /getRecommendatoon
func buy(w http.ResponseWriter, req *http.Request) {
	Logger.Info("/buy")
	defer Logger.Info("returned ", http.StatusAccepted)
	// Get user (session) identifier, for example by inspection of header X-User
	user := req.Header["X-User"][0]

	// export metric to metrics database; this is best effort; ignore any failure
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, _ = (*client).WriteMetric(
		ctx,
		&pb.MetricValue{
			Name:        "sample_metric",
			Value:       fmt.Sprintf("%f", rand.Float64()*100.0), // strconv.Itoa(rand.Intn(100)),
			Application: "default/backend",
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

	// establish connection to ABn service
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", getAbnService(), getAbnServicePort()), opts...)
	if err != nil {
		panic("Cannot establish connection with abn service")
		// return
	}
	c := pb.NewABNClient(conn)
	client = &c

	// configure frontend service with "/hello" and "/goodbye" endpoints
	http.HandleFunc("/getRecommendation", getRecommendation)
	http.HandleFunc("/buy", buy)
	http.ListenAndServe(":8090", nil)
}

func getAbnService() string {
	if value, ok := os.LookupEnv("ABN_SERVICE"); ok {
		return value
	}
	return "iter8"
}

func getAbnServicePort() string {
	if value, ok := os.LookupEnv("ABN_SERVICE_PORT"); ok {
		return value
	}
	return "50051"
}
