.PHONY: setup run db dbmigrate

APP_NAME := S3_FriendManagementAPI_NhutTo
APP_PATH := /${APP_NAME}
COMPOSE := docker-compose -f docker-compose.yaml
RUN_COMPOSE = $(COMPOSE) run --rm --service-ports -w $(APP_PATH) $(MOUNT_VOLUME) go

setup: db sleep dbmigrate

run: 
	@$(RUN_COMPOSE) env $(shell cat .env.dev | egrep -v '^#|^DATABASE_URL' | xargs) \
		go run -mod=vendor main.go

test: 
	@$(RUN_COMPOSE) env $(shell cat .env.dev | egrep -v '^#|^DATABASE_URL' | xargs) \
		go test -mod=vendor ./... -v

db:
	$(COMPOSE) up -d db

dbmigrate: MOUNT_VOLUME = $(if $(strip $(CONTAINER_SUFFIX)),,-v $(shell pwd)/db/migrations:/migrations)
dbmigrate:
	$(COMPOSE) run --rm $(MOUNT_VOLUME) db-migrate \
	sh -c './migrate -path /migrations -database $$DATABASE_URL up'

build:
	go clean -mod=vendor -i -x -cache ./...
	go build -mod=vendor -v -a -i ./

teardown:
	$(COMPOSE) down -v
	$(COMPOSE) rm --force --stop -v

sleep:
	sleep 10

graph-gen:
	go run github.com/99designs/gqlgen generate

update-vendor:
	$(RUN_COMPOSE) sh -c "go mod tidy && go mod vendor"