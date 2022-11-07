// GENERATED CODE -- DO NOT EDIT!

// Original file comments:
// protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative     abn/grpc/abn.proto
// python -m grpc_tools.protoc -I../../../iter8-tools/iter8/abn/grpc --python_out=. --grpc_python_out=. ../../../iter8-tools/iter8/abn/grpc/abn.proto 
//
'use strict';
var grpc = require('@grpc/grpc-js');
var abn_pb = require('./abn_pb.js');
var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js');

function serialize_google_protobuf_Empty(arg) {
  if (!(arg instanceof google_protobuf_empty_pb.Empty)) {
    throw new Error('Expected argument of type google.protobuf.Empty');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_google_protobuf_Empty(buffer_arg) {
  return google_protobuf_empty_pb.Empty.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_main_Application(arg) {
  if (!(arg instanceof abn_pb.Application)) {
    throw new Error('Expected argument of type main.Application');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_main_Application(buffer_arg) {
  return abn_pb.Application.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_main_MetricValue(arg) {
  if (!(arg instanceof abn_pb.MetricValue)) {
    throw new Error('Expected argument of type main.MetricValue');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_main_MetricValue(buffer_arg) {
  return abn_pb.MetricValue.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_main_Session(arg) {
  if (!(arg instanceof abn_pb.Session)) {
    throw new Error('Expected argument of type main.Session');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_main_Session(buffer_arg) {
  return abn_pb.Session.deserializeBinary(new Uint8Array(buffer_arg));
}


// for more information, see https://github.com/iter8-tools/iter8/issues/1257
//
var ABNService = exports.ABNService = {
  // Identify a version the caller should send a request to.
// Should be called for each request (transaction).
lookup: {
    path: '/main.ABN/Lookup',
    requestStream: false,
    responseStream: false,
    requestType: abn_pb.Application,
    responseType: abn_pb.Session,
    requestSerialize: serialize_main_Application,
    requestDeserialize: deserialize_main_Application,
    responseSerialize: serialize_main_Session,
    responseDeserialize: deserialize_main_Session,
  },
  // Write a metric value to metrics database.
// The metric value is explicitly associated with a list of transactions that contributed to its computation.
// The user is expected to identify these transactions.
writeMetric: {
    path: '/main.ABN/WriteMetric',
    requestStream: false,
    responseStream: false,
    requestType: abn_pb.MetricValue,
    responseType: google_protobuf_empty_pb.Empty,
    requestSerialize: serialize_main_MetricValue,
    requestDeserialize: deserialize_main_MetricValue,
    responseSerialize: serialize_google_protobuf_Empty,
    responseDeserialize: deserialize_google_protobuf_Empty,
  },
};

exports.ABNClient = grpc.makeGenericClientConstructor(ABNService);
