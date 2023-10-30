#!/bin/bash

# constants
SERVICE="http://localhost:8090"
GET_RECOMMENDATION="$SERVICE/getRecommendation"
BUY="$SERVICE/buy"

while (( 1 )); do 
    __user=$(uuidgen)
    __num_purchases=$(( ( RANDOM % 5 ) +1 ))
    j=1
    while (( ${j} <= ${__num_purchases} )); do
        echo "purchase $j of $__num_purchases for user $__user"
        __num_recommendations=$(( ( RANDOM % 5 )  + 1 ))
        # get some recommendations
        i=1
        while (( ${i} <= ${__num_recommendations} )); do
            echo "> recommendation $i of $__num_recommendations"
            curl -s ${GET_RECOMMENDATION} -H "X-User: ${__user}"
            sleep $(( ( RANDOM % 2000 ) / 1000 ))
            (( i += 1 ))
        done

        # buy
        curl -s ${BUY} -H "X-User: ${__user}"
        echo

        (( j += 1 ))
    done
    echo

    sleep 1
done