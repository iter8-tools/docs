var grpc = require('@grpc/grpc-js');
var messages = require('./abn_pb.js');
var services = require('./abn_grpc_pb.js');
var http = require('http');
var random = require('random-number');

'use strict';
const express = require('express');
const { registerChannelzSubchannel } = require('@grpc/grpc-js/build/src/channelz.js');
const { application } = require('express');
const { getLogger } = require('@grpc/grpc-js/build/src/logging.js');

const app  = express();

// define map of track to route to backend service
const trackToRoute = [
    "http://backend.default.svc.cluster.local:8091",
    "http://backend-candidate-1.default.svc.cluster.local:8091",
]

// establish connection to ABn service
var abnService = process.env.ABN_SERVICE || 'iter8'
var abnServicePort = process.env.ABN_SERVICE_PORT || 50051
var abnEndpoint = abnService + ':' + abnServicePort.toString()
var client = new services.ABNClient(abnEndpoint, grpc.credentials.createInsecure());

// /getRecommendation endpoint; calls backend service /recommend endpoint
app.get('/getRecommendation', (req, res) => {
    console.info('/getRecommendation')

    // identify default route
    route = trackToRoute[0];

    // call ABn service API Lookup() to get an assigned track for the user
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
            track = Number(session.getTrack())
            if (track != NaN && 0 <= track && track < trackToRoute.length) {
                route = trackToRoute[track]
            }
        }

        console.info('lookup suggested route %s', route)

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
            console.error("ERROR: " + err.syscall + " " + err.code + " " + err.hostname);
            res.status(500).send(err.message);
        });
    });
});

// /buy endpoint
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
    res.send(`Purchase complete`)
});

app.listen(8090, '0.0.0.0');
