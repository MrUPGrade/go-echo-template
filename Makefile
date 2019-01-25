GO:=go
DC:=docker-compose
DOCKER:=docker

DATE=$(shell date)



build-app:
	@-echo "### building echoapi [$(DATE)"
	$(GO) build -o echoapi

build-docker:
	@-echo "### building echoapi docker [$(DATE)]"
	$(DOCKER) build -t mrupgrade/echoapi:latest .

build: build-app build-docker


app-docker-up:
	$(DOCKER) run -d --name echoapi -p 8080:8080 mrupgrade/echoapi:latest

app-docker-down:
	-$(DOCKER) rm -f echoapi

app-docker-setup: app-docker-down app-docker-up


env-dev-up:
	$(DC) up -d

env-dev-down:
	$(DC) down

env-dev-setup: env-dev-down env-dev-up

test-run: build app-docker-setup
