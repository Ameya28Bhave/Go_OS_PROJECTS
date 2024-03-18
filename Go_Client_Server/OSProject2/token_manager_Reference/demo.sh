#!/bin/bash

PORT=$1
HOST=localhost

# output directories and the files
mkdir -p output
rm -rf output/server_op.txt
touch output/server_op.txt

# launching the server
echo "Lauching the server"
set -x
go run server.go -port $PORT > output/server_op.txt 2>&1 &
{ set +x; } 2>/dev/null
SERVER_PID=$!


echo "Server Port no.: $PORT"
echo "Server PID is: $SERVER_PID"
sleep 4
echo "We have launched the sever"

# launching the client requests
echo "Launching the client requests"
j=1
while [ $j -le 10 ]
do
    rm -rf output/client_$j.txt
    touch output/client_$j.txt
    echo "Creating the request: id=$j" > output/client_$j.txt
    set -x
    go run client.go -create -id $j -host $HOST -port $PORT >> output/client_$j.txt 2>&1 &
    { set +x; } 2>/dev/null
    sleep 0.4
    echo "Write Request: id=$j, -name=a$j, low=0, mid=100000, high=5000000" >> output/client_$j.txt
    set -x
    go run client.go -write -id $j -name a$j -low 0 -mid 100000 -high 5000000 -host $HOST -port $PORT >> output/client_$j.txt 2>&1 &
    { set +x; } 2>/dev/null
    sleep 0.4
    echo "Read the request: id=$j" >> output/client_$j.txt
    set -x
    go run client.go -read -id $j -host $HOST -port $PORT >> output/client_$j.txt 2>&1 &
    { set +x; } 2>/dev/null
    sleep 0.4
    echo "Drop the request: id=$j" >> output/client_$j.txt
    set -x
    go run client.go -drop -id $j -host $HOST -port $PORT >> output/client_$j.txt 2>&1 &
    { set +x; } 2>/dev/null
    let j++
done
echo "Launched the client requests"

# wait till all request got servered
k=1
requests=$(( 4*j-4 ))
while true
do
    PROCESSED_REQ="$(grep "Processed request number" output/server_op.txt | tail -n 1 | cut -w -f 4)"
    if [[ $PROCESSED_REQ =~ $( printf '%d' $requests ) || $k == 30 ]]; then
        echo "All requests either processed or timeout reached"
        break
    else
        echo "Waiting for all requests getting processed..."
        sleep 4
        let k++
    fi
done

echo "Closing the server"
sleep 1
pkill -9 -P $SERVER_PID
echo "Server is closed"

echo "**** Please check outputs in 'output' directory ****"