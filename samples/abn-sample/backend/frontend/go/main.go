package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	pb "github.com/kalantar/ab-example/frontend/go/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// var log *logrus.Logger

var (
	// map of track to route to backend service
	trackToRoute = map[string]string{
		"default":   "http://backend:8091",
		"candidate": "http://backend-candidate:8091",
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
	route := trackToRoute["default"]

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
	// if successful, use recommended track; otherwise will use default route
	if err == nil && s != nil {
		r, ok := trackToRoute[s.GetTrack()]
		if ok {
			route = r
		}
	}

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
	w.WriteHeader(http.StatusAccepted)
}

func main() {
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		DisableQuote:    true,
		DisableSorting:  true,
	})
	Logger.SetLevel(logrus.TraceLevel)

	t := time.Now()
	fmt.Println(t)
	eu := t.Unix()
	fmt.Println(eu)
	ef := float64(eu)
	fmt.Println(ef)
	// time.Unix(int64(math.Round(encoded[5])), 0)
	di := int64(math.Round(ef))
	fmt.Println(di)
	dt := time.Unix(di, 0)
	fmt.Println(dt)

	// establish connection to ABn service
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial("iter8-abn:50051", opts...)
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
