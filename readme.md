MQTT:
- [mqtt.org](https://mqtt.org)  
- [hivemq.com](https://www.hivemq.com)
- [Video tutorial with theory](https://www.youtube.com/playlist?list=PLRkdoPznE1EMXLW6XoYLGd4uUaB6wB0wd)

Soft:
Most commonly used libraries are part of [Eclipse Paho](https://eclipse.dev/paho/index.php?page=downloads.php) project
- [Clients](https://eclipse.dev/paho/index.php?page=downloads.php)
- [Broker](https://mosquitto.org/)

Repo contains 2 test clients (publisher and subscriber) written with Go.  
Publisher (Sensor) is a mock of Temperature and Humidity sensor which periodically sends sensor data to broker.  
Subscriber (Controller) is a mock of some controlling device.

Mosquitto is run from [docker container](https://hub.docker.com/_/eclipse-mosquitto)  
- To pull image: `docker pull eclipse-mosquitto`
- To create and run container: `start_mosquitto.sh`
- To run existing container: `docker start -a mosquitto`  

