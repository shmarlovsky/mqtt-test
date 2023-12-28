#!/bin/bash

# named container is created, next run should be done wiht `docker start` command
docker run -it -p 1883:1883 -p 9001:9001 -v $(pwd)/mosquitto.conf:/mosquitto/config/mosquitto.conf -v /mosquitto/data -v /mosquitto/log --name mosquitto eclipse-mosquitto
