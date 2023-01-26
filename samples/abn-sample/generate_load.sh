#!/bin/bash

# constants
SERVICE="http://localhost:8090"
GET_RECOMMENDATION="$SERVICE/getRecommendation"
BUY="$SERVICE/buy"

while (( 1 )); do 
    __user=$(uuidgen)
    __num_recommendations=$(( ( RANDOM % 5 )  + 1 ))
    # get some recommendations
    i=0
    while (( ${i} < ${__num_recommendations} )); do
        curl -s ${GET_RECOMMENDATION} -H "X-User: ${__user}"
        sleep $(( ( RANDOM % 2000 ) / 1000 ))
        (( i += 1 ))
    done

    # buy
    curl -s ${BUY} -H "X-User: ${__user}"
    echo

    sleep 1
done