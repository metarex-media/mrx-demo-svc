#!/usr/bin/env bash

# run the api
cd api || exit
# run it in the background
./api &

# run the main server
cd .. && ./mrx-demo-svc
