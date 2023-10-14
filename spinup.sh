#!/bin/bash

# only for testing

while getopts ":p:" opt; do
    case $opt in
        p) PORT="$OPTARG"
        ;;
        \?) echo "Invalid arguments"
        exit 1
        ;;
    esac
done

[ -z $PORT ] && echo "No port set! set port with -p flags" && exit 1

go run main.go &

while true; do
    inotifywait -e modify -q ./**.go;
    pkill -P $$ --signal SIGINT
    fuser --silent -k -n tcp $PORT
    go run main.go &
done

