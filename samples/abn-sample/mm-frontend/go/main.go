package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/kalantar/ab-example/mm-frontend/go/inference"
	"github.com/sirupsen/logrus"

	abn "github.com/iter8-tools/iter8/abn/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var Logger = logrus.New()

// map of version number to route to backend service
// here route is a model id of the model to which the inference request should be sent
var versionNumberToRoute = []string{
	"backend-0",
	"backend-1",
}

// implment /getRecommendation endpoint
// calls backend service (ML model served in modelmesh-serving)
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

	// call ABn service API Lookup() to get a recommended version for the user
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s, err := getABnClient().Lookup(
		ctx,
		&abn.Application{
			Name: backendName(),
			User: user,
		},
	)
	// if successful, use recommended version; otherwise will use default
	if err != nil {
		Logger.Info("error: " + err.Error())
	}
	// if successful, use returned version; otherwise will use the default
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
			Application: "modelmesh-serving/backend",
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
		// establish connection to ABn service
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
	return lookupEnv("BACKEND_APPLICATION_NAME", "modelmesh-serving/backend")
}

// callBackend calls inference service with mm-vmodel-id=$route using modelmesh-serving gRPC API
// equivalent to:
//
//	grpcurl -plaintext -proto proto -d data \
//	   -H 'mm-vmodel-id: $route' \
//	   modelmesh-serving.modelmesh-serving:8033 inference.GRPCInferenceService.ModelInfer
//
// input data is hard-coded in this example
func callBackend(route string) (string, error) {
	Logger.Infof("callBackend (%s)", route)
	defer Logger.Info("callBackend finished")

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("mm-vmodel-id", route))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// send request
	resp, err := getBackendClient().ModelInfer(
		ctx,
		&inference.ModelInferRequest{
			Inputs: []*inference.ModelInferRequest_InferInputTensor{
				{
					Name:     "predict",
					Shape:    []int64{1, 64},
					Datatype: "FP32",
					Contents: &inference.InferTensorContents{
						Fp32Contents: []float32{0.0, 0.0, 1.0, 11.0, 14.0, 15.0, 3.0, 0.0, 0.0, 1.0, 13.0, 16.0, 12.0, 16.0, 8.0, 0.0, 0.0, 8.0, 16.0, 4.0, 6.0, 16.0, 5.0, 0.0, 0.0, 5.0, 15.0, 11.0, 13.0, 14.0, 0.0, 0.0, 0.0, 0.0, 2.0, 12.0, 16.0, 13.0, 0.0, 0.0, 0.0, 0.0, 0.0, 13.0, 16.0, 16.0, 6.0, 0.0, 0.0, 0.0, 0.0, 16.0, 16.0, 16.0, 7.0, 0.0, 0.0, 0.0, 0.0, 11.0, 13.0, 12.0, 1.0, 0.0},
					},
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.GetModelName(), err
}

var backendClient *inference.GRPCInferenceServiceClient

func getBackendClient() inference.GRPCInferenceServiceClient {
	if backendClient == nil {
		// establish connection to backend ML service
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		conn, err := grpc.Dial(
			fmt.Sprintf(
				"%s:%s",
				lookupEnv("RECOMMENDATION_SERVICE", "modelmesh-serving"),
				lookupEnv("RECOMMENDATION_SERVICE_PORT", "8033"),
			),
			opts...,
		)
		if err != nil {
			panic("Cannot establish connection with ML service")
			// return
		}
		c := inference.NewGRPCInferenceServiceClient(conn)
		backendClient = &c
	}

	return *backendClient
}
