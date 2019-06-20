SHELL=/bin/bash
BINNAME:=echoapi

GO:=go
DOCKER:=docker
DOCKER_IMAGE:=mrupgrade/$(BINNAME):latest

DATE=$(shell date)

DEV_ENV= source dev-env.sh;
DC=$(DEV_ENV) docker-compose -p $(BINNAME)



build-app:
	@-echo "### building $(BINNAME) [$(DATE)"
	$(GO) build -o $(BINNAME) cmd/api/main.go

build-docker:
	@-echo "### building $(BINNAME) docker [$(DATE)]"
	$(DOCKER) build -f build/Dockerfile -t mrupgrade/$(BINNAME):latest .

build: build-app build-docker



publish-docker:
	$(DOCKER) push $(DOCKER_IMAGE)



app-docker-up:
	$(DOCKER) run -d --name $(BINNAME) --env-file .dev.env -p 8080:8080 $(DOCKER_IMAGE)

app-docker-down:
	-$(DOCKER) rm -f $(BINNAME)

app-docker-setup: app-docker-down build app-docker-up



env-dev-up:
	$(DC) up -d

env-dev-down:
	$(DC) down

env-dev-setup: env-dev-down env-dev-up



test-run: build app-docker-setup
