SHELL=/bin/bash
GO:=go
DOCKER:=docker

DATE=$(shell date)
DEV_ENV=. dev-env.sh;
DC=$(DEV_ENV) docker-compose


build-app:
	@-echo "### building echoapi [$(DATE)"
	$(GO) build -o echoapi

build-docker:
	@-echo "### building echoapi docker [$(DATE)]"
	$(DOCKER) build -t mrupgrade/echoapi:latest .

build: build-app build-docker


app-docker-up:
	$(DOCKER) run -d --name echoapi --env-file -p 8080:8080 mrupgrade/echoapi:latest

app-docker-down:
	-$(DOCKER) rm -f echoapi

app-docker-setup: app-docker-down app-docker-up


env-dev-up:
	$(DC) up -d

env-dev-down:
	$(DC) down

env-dev-setup: env-dev-down env-dev-up

test-run: build app-docker-setup
