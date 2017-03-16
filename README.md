# fwatcher
An application to watch a given directory for new files, read it and publish to Kafka ( using actors )

### Install:
    1. Run a temporary Kafka-Zookeeper cluster: `docker-compose up -d`. The docker-compose file is in the root of project.
    2. go get -u github.com/ansrivas/fwatcher
    3. This will start monitoring the current dir: `fwatcher --config config.yaml`
    4. Just copy any random file in here. 
