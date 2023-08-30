import random
import logging

from http import HTTPStatus
import os
import requests

import grpc
import abn_pb2
import abn_pb2_grpc

from flask import Flask, request

# Map of version number to route to backend service
versionNumberToRoute = [
    "http://backend.default.svc.cluster.local:8091",
    "http://backend-candidate-1.default.svc.cluster.local:8091"
]

app = Flask(__name__)
app.logger.setLevel(logging.DEBUG)

# Implement /getRecommendation endpoint
# Calls backend service /recommend endpoint
@app.route('/getRecommendation')
def getRecommendation():
    # Get user (session) identifier, for example, by inspection of header X-User
    user = request.headers['X-User']

    # Get endpoint of backend endpoint "/recommend"
    # In this example, the backend endpoint depends on the version of the backend service
    # The user is assigned by the Iter8 SDK Lookup() method

    # Start with default route
    route = versionNumberToRoute[0]

    # Establish connection to ABn service
    abnSvc = os.getenv('ABN_SERVICE', 'iter8') + ":" + os.getenv('ABN_SERVICE_PORT', '50051')
    with grpc.insecure_channel(abnSvc) as channel:
        stub = abn_pb2_grpc.ABNStub(channel)

        try:
            # Call A/B/n service API Lookup() to get a recommended version number for the user
            s = stub.Lookup( \
                abn_pb2.Application(name="default/backend", \
                user=user) \
            )

            # Lookup route using version number
            route = versionNumberToRoute[int(s.versionNumber)]
        except Exception as e:
            # Use default
            app.logger.error("error: %s", e)
            pass

    app.logger.info('lookup suggested route %s', route)

    # Call backend service using URL
    try:
        r = requests.get(url=route + "/recommend", allow_redirects=True)
        r.raise_for_status()
        recommendation = r.text
    except Exception as e:
        return "call to backend endpoint /recommend failed: {0}".format(e), HTTPStatus.INTERNAL_SERVER_ERROR

    return "Recommendation: {0}".format(recommendation)
    
# Implement /buy endpoint
# Writes value for sample_metrc which may have spanned several calls to /getRecommendation
@app.route('/buy')
def buy():
    # Get user (session) identifier, for example, by inspection of header X-User
    user = request.headers['X-User']

	# Export metric to metrics database
	# This is best effort; we ignore any failure

    # establish connection to ABn service
    abnSvc = os.getenv('ABN_SERVICE', 'iter8') + ":" + os.getenv('ABN_SERVICE_PORT', '50051')
    with grpc.insecure_channel(abnSvc) as channel:
        stub = abn_pb2_grpc.ABNStub(channel)

        # Export metric to metrics database
        # This is best effort; we ignore any failure
        try:
            stub.WriteMetric( \
                abn_pb2.MetricValue(name="sample_metric", \
                value=str(random.randint(0,100)), \
                application="default/backend", \
                user=user) \
            )
        except Exception as e:
            pass

    return "Purchase complete"
