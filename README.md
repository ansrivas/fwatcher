# fwatcher

An application to watch a given directory for new files, read it and publish to Kafka ( using actors )

## Install:

1. Run a temporary Kafka-Zookeeper cluster: `docker-compose up -d`. The docker-compose file is in the root of project.
2. `go get -u github.com/ansrivas/fwatcher`
3. This will start monitoring the current dir: `fwatcher --config config.yaml`
4. Just copy any random file in here.

## Usage:

--------------------------------------------------------------------------------

```bash
$ make
help:           Show available options with this Makefile
test:           Run all the tests
clean:          Clean the application and remove all the docker containers.
build:          Build the application
app_help:       Display flags accepted by the application
test_run:       Run the application in a test mode
_recreate_env:  Recreate the docker environment and create a default topic.
migrate:        Run migration to populate the db
dock_run_fg:    Run docker containers, foreground.
dock_run_bg:    Run docker containers, background.
build_docker:   Build docker containers
sys_info:       Show docker containers info
```
