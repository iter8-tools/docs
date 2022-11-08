import random

from http import HTTPStatus
import os
import requests

import grpc
import abn_pb2
import abn_pb2_grpc

from flask import Flask, request

# map of track to route to backend service
trackToRoute = {
    "default": "http://backend:8091",
	"candidate": "http://backend-candidate:8091"
}

app = Flask(__name__)

# implement /getRecommendation endpoint
# calls backend service /recommend endpoint
@app.route('/getRecommendation')
def getRecommendation():
    # Get user (session) identifier, for example, by inspection of header X-User
    user = request.headers['X-User']

    # Get endpoint of backend endpoint "/recommend"
    # In this example, the backend endpoint depends on the version (track) of the backend service
    # the user is assigned by the Iter8 SDK Lookup() method

    # start with default route
    route = trackToRoute["default"]

    # establish connection to ABn service
    abnSvc = os.getenv('ABN_SERVICE', 'iter8-abn') + ":" + os.getenv('ABN_SERVICE_PORT', '50051')
    with grpc.insecure_channel(abnSvc) as channel:
        stub = abn_pb2_grpc.ABNStub(channel)

        try:
            # call ABn service API Lookup() to get an assigned track for the user
            s = stub.Lookup( \
                abn_pb2.Application(name="default/backend", \
                user=user) \
            )

            # lookup route using track
            route = trackToRoute[s.track]
        except Exception as e:
            # use default
            pass

    # call backend service using url
    try:
        r = requests.get(url=route + "/recommend", allow_redirects=True)
        r.raise_for_status()
        recommendation = r.text
    except Exception as e:
        return "call to backend endpoint /recommend failed: {0}".format(e), HTTPStatus.INTERNAL_SERVER_ERROR

    return "Recommendation: {0}".format(recommendation)
    
# implement /buy endpoint
# writes value for sample_metrc which may have spanned several calls to /getRecommendation
@app.route('/buy')
def buy():
    # Get user (session) identifier, for example, by inspection of header X-User
    user = request.headers['X-User']

	# export metric to metrics database
	# this is best effort; we ignore any failure

    # establish connection to ABn service
    abnSvc = os.getenv('ABN_SERVICE', 'iter8-abn') + ":" + os.getenv('ABN_SERVICE_PORT', '50051')
    with grpc.insecure_channel(abnSvc) as channel:
        stub = abn_pb2_grpc.ABNStub(channel)

        # export metric to metrics database
        # this is best effort; we ignore any failure
        try:
            stub.WriteMetric( \
                abn_pb2.MetricValue(name="sample_metric", \
                value=str(random.randint(0,100)), \
                application="default/backend", \
                user=user) \
            )
        except Exception as e:
            pass

    return ""