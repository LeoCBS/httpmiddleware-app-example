version ?= latest
imgdev = leocbs/httpmiddleware-app-example:$(version)
RUNCOMPOSE=docker-compose run --rm httpmiddleware-app-example


imagedev:
	docker build --target devimage . -t $(imgdev)

check-integration: imagedev
	$(RUNCOMPOSE) go test -tags=integration -timeout 60s -race ./...

static-analysis: imagedev
	$(RUNCOMPOSE) golangci-lint run ./...

modtidy:
	go mod tidy

run-compose: imagedev
	docker-compose run --rm --service-ports --entrypoint "go run main.go" httpmiddleware-app-example

