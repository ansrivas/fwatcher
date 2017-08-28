DEP := $(shell command -v dep 2>/dev/null)
BUILD_TIME=`date -u +%FT%T%z`

.DEFAULT_GOAL := help

help:          ## Show available options with this Makefile
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY : test
test:          ## Run all the tests
test:
	./test.sh

dep:
	go get -u github.com/golang/dep/cmd/dep
ifndef DEP
	make dep
endif
	dep ensure

clean:         ## Clean the application and remove all the docker containers.
	@go clean -i ./...
	@rm -rf ./fwatcher
	@docker-compose down -v

build:         ## Build the application
build:	clean
	@go build github.com/ansrivas/fwatcher

.PHONY : app_help
app_help:      ## Display flags accepted by the application
APP_HELP = "$(shell ./fwatcher --help)"
app_help: build
	@echo $(APP_HELP)

.PHONY: test_run
test_run:      ## Run the application in a test mode
test_run:	_recreate_env build
	@echo "Running now.."
	@./fwatcher --config ./config.yaml

_recreate_env: ## Recreate the docker environment and create a default topic.
_recreate_env:	clean
	docker-compose up -d && \
	chmod +x ./wait-for-it.sh && \
	./wait-for-it.sh localhost:19092 --timeout=0 --	docker exec -it kafka-01-c /usr/bin/kafka-topics --create --zookeeper localhost:22181 --replication-factor 1 --partitions 100 --topic test_topic

migrate:       ## Run migration to populate the db
migrate:
	@go get -u github.com/pressly/goose/cmd/goose
	@goose -dir db/migrations/ postgres "user=testuser password=testpassword123 dbname=testdb sslmode=disable" $(command)

.PHONY: dock_run_fg
dock_run_fg:   ## Run docker containers, foreground.
dock_run_fg:	build_docker
	docker-compose up

.PHONY: dock_run_bg
dock_run_bg:   ## Run docker containers, background.
dock_run_bg:
	docker-compose up -d

.PHONY: build_docker
build_docker:  ## Build docker containers
build_docker:
	docker-compose build

.PHONY: sys_info
sys_info:      ## Show docker containers info
sys_info:
	docker system df -v
