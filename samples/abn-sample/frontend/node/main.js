var messages = require('./abn_pb.js');
var services = require('./abn_grpc_pb.js');
var grpc = require('@grpc/grpc-js');
var http = require('http');
var random = require('random-number');

'use strict';
const express = require('express');
const { registerChannelzSubchannel } = require('@grpc/grpc-js/build/src/channelz.js');
const { application } = require('express');

const app  = express();

// map of track to route to backend service
const trackToRoute = {
    "default":   "http://backend:8091",
    "candidate": "http://backend-candidate:8091",
}
// establish connection to ABn service
var abnService = process.env.ABN_SERVICE || 'iter8-abn'
var abnServicePort = process.env.ABN_SERVICE_PORT || 50051
var abnEndpoint = abnService + ':' + abnServicePort.toString()
console.log(abnEndpoint)
var client = new services.ABNClient(abnEndpoint, grpc.credentials.createInsecure());

// implement /getRecommendation endpoint
// calls backend service /recommend endpoint
app.get('/getRecommendation', (req, res) => {
	// Get user (session) identifier, for example by inspection of header X-User
    const user = req.header('X-User')

    // Get endpoint of backend endpoint "/recommend"
	// In this example, the backend endpoint depends on the version (track) of the backend service
	// the user is assigned by the Iter8 SDK Lookup() method

    // start with default route
    route = trackToRoute['default']

    // call ABn service API Lookup() to get an assigned track for the user
    var application = new messages.Application();
    application.setName('default/backend');
    application.setUser(user);
    client.lookup(application, function(err, session) {
        // lookup route using track
        route = trackToRoute[session.getTrack()];

        // call backend service using route
        http.get(route + '/recommend', (resp) => {
            let str = '';
            resp.on('data', function(chunk) {
                str += chunk;
            });
            resp.on('end', function () {
                // write response to query
                res.send(`Recommendation: ${str}`);
            });
        }).on("error", (err) => {
            res.status(500).send(err.message);
        });
    });
});

// implement /buy endpoint
// writes value for sample_metric which may have spanned several calls to /getRecommendation
app.get('/buy', (req, res) => {
	// Get user (session) identifier, for example by inspection of header X-User
    const user = req.header('X-User')

	// export metric to metrics database
	// this is best effort; we ignore any failure

    // export metric
    var mv = new messages.MetricValue();
    mv.setName('sample_metric');
    mv.setValue(random({min: 0, max: 100, integer: true}).toString()); 
    mv.setApplication('default/backend');
    mv.setUser(user);
    client.writeMetric(mv, function(err, session) {});
    res.sendStatus(200);
});

app.listen(8090, '0.0.0.0');
