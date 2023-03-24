---
template: main.html
---

# A/B/n Testing

A/B/n testing relies on business metrics typically computed by a frontend, user-facing, application component. 

![A/B/n experiment](images/abn.png)

Metric values often depend on one or more interactions with backend (not user-facing) application components. To run an A/B/n test on a backend component, it is necessary to be able to associate a metric value (computed by the frontend component) to the version of the backend component that contributed to its computation. 
The challenge is that the frontend component often does not know which version of the backend component processed a given request. To address this challenge, Iter8 introduces an A/B/n SDK. 

The Iter8 SDK introduces the concept of a *track identifier*. A track is a logical version of a Kubernetes application. The set of valid track identifiers is fixed over the lifetime of the application. The version of the application associated with a given track identifier changes over time as new versions are developed. For a given application, the set of track identifiers is fixed; the number of track identifiers determines how many versions of the application can be deployed/tested at the same time. Because the set of track identifiers is fixed, they can be used to configure routing to the application.

The Iter8 SDK provides two APIs to frontend application components:

a. **Lookup()** - Given an application and user session, returns a track identifier. So long as there are no changes in configuration, the track identifier (and hence the route) will be same for the same user session, guaranteeing session stickiness.

b. **WriteMetric()** -  Given an application, a user session, a metric name its value, *WriteMetric()* associates the metric value with the appropriate version of the application. 

## Configuring the Iter8 A/B/n Service

An Iter8 A/B/n service implements the gRPC API. This service is configured, at deployment, to watch the resource objects for a set of applications so that it can identify new versions and their mapping to a track identifier.

To watch for versions of an application, specify the list of the types of the objects that must be present and ready for a version to be considered ready:

`--set "apps.<namespace>.<application_name>.resources={<comma separated list resoure types>}"`

For example, to watch for versions of an application `my_app` in the namespace `my_namespace` where each version is composed of a Kubernetes service object and a Kubernetes deployment object, specify:

`--set "apps.my_namespace.my_app.resources={service,deployment}"`

Valid resource types are corresponding Kubernetes resource types (specified by group, version and resource) are listed below. When the required condition value is `true`, the resource object is considered ready.

| Type Name | Kubernetes Resource Type (GVR) | Required Condition |
| ---- | ---- | ----------- |
| service  | v1 services | - |
| deployment | apps/v1 deployments | Available |
| ksvc | serving.knative.dev/v1 services | Ready |

If more than one candidate version can be deployed at the same time, specify the maximum number using `maxNumCandidates`:

`--set apps.<namespace>.<application_name>.maxNumCandidates=<number>`

From the above configuration, Iter8 infers the names of the expected resource objects using these assumptions:

- The baseline track identifier is the application name
- Track identifiers associated with candidate versions are of the form `<application_name>-candidate-<index>`
- All resource objects for all versions are deployed in the same namespace
- There is only 1 resource object of a given type in each version
- The name of each object in the version associated with the baseline track is the application name
- The name of each object in the version associate with a candidate track is of the form  `<application_name>-candidate-<index>` where index is 1, 2, etc.

## Deployment Time Configuration of Backend Components

As versions of the backend component are deployed or deleted, the Iter8 A/B/n service maintains a mapping of track identifier to available version. Using this mapping it is then able to respond appropriately to *Lookup()* and *WriteMetric()* requests.

To build and maintain it's mapping, the A/B/n service watches the resource objects specified as part of the A/B/n service configuration (see above). 
In particular, the configuration requires that the Kubernetes objects comprising the backend component adhere to the specified naming convention. Further, they should have the label `app.kubernetes.io/version` set to the version identifier.

## Developing Frontend Components: Using the SDK

The basic steps to author a frontend application component using the Iter8 SDK are outlined below for *Node.js* and *Go*. Similar steps would be required for any gRPC supported language.

### Use/Import language specific libraries

The gRPC protocol buffer definition is used to generate language specific implementation. These files can be used directly or packaged and imported as a library. As examples, the [Node.js sample](https://github.com/iter8-tools/docs/tree/main/samples/abn-sample/frontend/node) uses manually generated files directly. On the other hand, the [Go sample](https://github.com/iter8-tools/docs/tree/main/samples/abn-sample/frontend/go) imports the library provided by the core Iter8 service implementation. In addition to the API specific methods, some general gRPC libraries are required.

=== "Node.js"
    The manually generated node files [`abn_pd.js`](https://raw.githubusercontent.com/iter8-tools/docs/main/samples/abn-sample/frontend/node/abn_pb.js) and [`abn_grpc_pb.js`](https://raw.githubusercontent.com/iter8-tools/docs/main/samples/abn-sample/frontend/node/abn_grpc_pb.js) used in the sample application can be copied and used without modification.

    ```javascript
    var grpc = require('@grpc/grpc-js');

    var messages = require('./abn_pb.js');
    var services = require('./abn_grpc_pb.js');
    ```

=== "Go"
    ```go
    import (
        "google.golang.org/grpc"
        "google.golang.org/grpc/credentials/insecure"

        pb "github.com/iter8-tools/iter8/abn/grpc"
    )
    ```

### Instantiate a gRPC client

Instantiate a client to the Iter8 A/B/n service:

=== "Node.js"
    ```javascript
    var client = new services.ABNClient(abnEndpoint, grpc.credentials.createInsecure());
    ```

=== "Go"
    ```go
    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
    conn, err := grpc.Dial(fmt.Sprintf("%s:%s", getAbnService(), getAbnServicePort()), opts...)
    if err != nil {
        panic("Cannot establish connection with abn service")
    }
    c := pb.NewABNClient(conn)
    client = &c
    ```

### Define routing

Track identifiers are mapped to a static set of endpoints. One approach is to maintain a map from track identifier to endpoint:

=== "Node.js"
    ```javascript
    const trackToRoute = {
        "backend":   "http://backend.default.svc.cluster.local:8091",
        "backend-candidate-1": "http://backend-candidate-1.default.svc.cluster.local:8091",
    }
    ```

=== "Go"
    ```go
    trackToRoute = map[string]string{
        "backend":             "http://backend.default.svc.cluster.local:8091",
        "backend-candidate-1": "http://backend-candidate-1.default.svc.cluster.local:8091",
    }
    ```

### Using *Lookup()*

Given a user session identifier, *Lookup()* returns a track identifier that can be used to route requests. In code sample below, the user session identifier is assumed to be passed in the `X-User` header of user requests. The track identifier is used as an index to the `trackToRoute` map defined above. A default is used if the call to *Lookup()* fails for any reason.

=== "Node.js"
    ```javascript
    var application = new messages.Application();
    application.setName('default/backend');
    application.setUser(req.header('X-User'));
    client.lookup(application, function(err, session) {
        if (err || (session.getTrack() == '')) {
            // use default route (see above)
            console.warn("error or null")
        } else {
            // use route determined by recommended track
            console.info('lookup suggested track %s', session.getTrack())
            route = trackToRoute[session.getTrack()];
        }

        // call backend service using route
        ...
    });
    ```

=== "Go"
    ```go
    route := trackToRoute["backend"]
    user := req.Header["X-User"][0]
    s, err := (*client).Lookup(
        ctx,
        &pb.Application{
            Name: "default/backend",
            User: user,
        },
    )
    if err == nil && s != nil {
        r, ok := trackToRoute[s.GetTrack()]
        if ok {
            route = r
        }
    }

    // call backend service using route
    ...
    ```

### Using *WriteMetric()*

As an example, a single metric named *sample_metric* is assigned a random value between 0 and 100 and written.

=== "Node.js"
    ```javascript
    var mv = new messages.MetricValue();
    mv.setName('sample_metric');
    mv.setValue(random({min: 0, max: 100, integer: true}).toString());
    mv.setApplication('default/backend');
    mv.setUser(user);
    ```

=== "Go"
    ```go
    _, _ = (*client).WriteMetric(
        ctx,
        &pb.MetricValue{
            Name:        "sample_metric",
            Value:       fmt.Sprintf("%f", rand.Float64()*100.0),
            Application: "default/backend",
            User:        user,
        },
    )
    ```
