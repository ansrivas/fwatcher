.DEFAULT_GOAL := help

help:          ## Show available options with this Makefile
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY : test
test:          ## Run all the tests
test:
	./test.sh

clean:         ## Clean the application
	@go clean -i ./...
	@rm -rf ./fwatcher

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
test_run:	clean build
	@echo "Running now.."
	@./fwatcher --config ./config.yaml

.PHONY: remove
remove:        ## Remove running docker containers and remove the volumes
remove:	clean
	docker-compose down -v

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
