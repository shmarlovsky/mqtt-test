Create test mqtt client with this behavior:
- client sends current temperature (random from interval) every N seconds
= in parallel client listen to incoming "commands" to send some additional data
- another client is kinda controller, gets data, sends commands
- local mosquitto broker
- depending on time, go and python clients


Mosquitto and it's docker:
https://mosquitto.org/download/
https://hub.docker.com/_/eclipse-mosquitto
`docker pull eclipse-mosquitto`
`docker run -it -p 1883:1883 -p 9001:9001 -v mosquitto.conf:/mosquitto/config/mosquitto.conf -v /mosquitto/data -v /mosquitto/log eclipse-mosquitto`
