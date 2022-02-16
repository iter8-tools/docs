---
template: main.html
---

# Benchmarking and Validating Streaming gRPC

!!! tip "Overview"
    This tutorial describes additional parameters enabled by `load-test-grpc` experiment while benchmarking and validating streaming gRPC.

***

???+ warning "Before you begin"
    1. [Install Iter8 CLI](../../getting-started/install.md).
    2. Get familiar with the [basic usage examples](usage.md).
    3. Run the RouteGuide service in a separate terminal. Choose any language and follow the linked instructions.

        === "C#"
            [Run the gRPC service](https://grpc.io/docs/languages/csharp/basics/#try-it-out).

        === "C++"
            [Run the gRPC service](https://grpc.io/docs/languages/cpp/basics/#try-it-out).

        === "Dart"
            [Run the gRPC service](https://grpc.io/docs/languages/dart/basics/#try-it-out).

        === "Go"
            [Run the gRPC service](https://grpc.io/docs/languages/go/basics/#try-it-out).

        === "Java"
            [Run the gRPC service](https://grpc.io/docs/languages/java/basics/#try-it-out).

        === "Kotlin"
            [Run the gRPC service](https://grpc.io/docs/languages/kotlin/basics/#try-it-out).

        === "Node"
            [Run the gRPC service](https://grpc.io/docs/languages/node/basics/#try-it-out).

        === "Objective-C"
            [Run the gRPC service](https://grpc.io/docs/languages/objective-c/basics/#try).

        === "PHP"
            [Run the gRPC service](https://grpc.io/docs/languages/php/basics/#try).

        === "Python"
            [Run the gRPC service](https://grpc.io/docs/languages/python/basics/#try-it-out).

        === "Ruby"
            [Run the gRPC service](https://grpc.io/docs/languages/ruby/basics/#try-it-out).

***

## Call data

For client streaming or bi-directional calls, the `load-test-grpc` chart accepts an array of messages, each element representing a single message within the stream call. For example, 

```shell
iter8 launch -c load-test-grpc \
          --set host="127.0.0.1:50051" \
          --set call="helloworld.Greeter.SayHello" \
          --set protoURL="https://raw.githubusercontent.com/grpc/grpc-go/master/examples/helloworld/helloworld/helloworld.proto" \
          --set data[0].name="Joe" \
          --set data[1].name="Kate" \
          --set data[2].name="Sara"
```
    
If a single object is given for data then it is automatically converted to an array with single element. In case of client streaming, `load-test-grpc` sends all the messages in the input array, and then closes and receives.

## Stream interval

Stream interval duration spreads stream sends by given amount, specified as a [Go duration string](https://pkg.go.dev/maze.io/x/duration#ParseDuration). This parameter applies to client and bidirectional streaming calls.

```shell
iter8 launch -c load-test-grpc \
          --set host="127.0.0.1:50051" \
          --set call="helloworld.Greeter.SayHello" \
          --set protoURL="https://raw.githubusercontent.com/grpc/grpc-go/master/examples/helloworld/helloworld/helloworld.proto" \
          --set data[0].name="Joe" \
          --set streamInterval="100ms"
```

***

## Stream call duration

This parameter sets the maximum stream call duration. For client streaming and bidirectional calls, `load-test-grpc` will send messages until this duration expires. For server streaming calls, `load-test-grpc` will receive messages until the duration has expired. In server streaming calls, expiration of this duration will result in a call cancelled error.

```shell
iter8 launch -c load-test-grpc \
          --set host="127.0.0.1:50051" \
          --set call="helloworld.Greeter.SayHello" \
          --set protoURL="https://raw.githubusercontent.com/grpc/grpc-go/master/examples/helloworld/helloworld/helloworld.proto" \
          --set data[0].name="Joe" \
          --set streamCallDuration="500ms"
```

***

## Stream call count

This parameter sets the maximum number of message sends or receives that will be performed by `load-test-grpc` in a streaming call, before closing the stream and ending the call. 

For client streaming and bidirectional calls, this represents the number of messages sent. For server streaming calls `load-test-grpc` will receive messages until the specified count is reached. Note that in server streaming calls, reaching this count will result in a call cancelled error.

```shell
iter8 launch -c load-test-grpc \
          --set host="127.0.0.1:50051" \
          --set call="helloworld.Greeter.SayHello" \
          --set protoURL="https://raw.githubusercontent.com/grpc/grpc-go/master/examples/helloworld/helloworld/helloworld.proto" \
          --set data[0].name="Joe" \
          --set streamCallCount=100
```

If the data array contains more elements than the count, only messages up to the specified count will be used. If the data array contains fewer elements than the count specified, the data will be iterated over until the specified count is reached.
