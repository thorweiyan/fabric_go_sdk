.PHONY: all dev clean build env-up env-down run

all: clean build env-up run

dev: build run

restart: env-down clean env-up

##### BUILD
build:
	@echo "Build ..."
	#@dep ensure
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@cd fixtures && docker-compose up --force-recreate -d
	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	@cd fixtures && docker-compose down
	@echo "Environment down"

##### RUN
run:
	@echo "Start app ..."
	@./fabric_go_sdk

##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@rm -rf /tmp/fudan-* fudan
	@docker rm -f -v `docker ps -a --no-trunc | grep "fudan" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "fudan" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"
