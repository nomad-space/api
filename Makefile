PROJECT?=github.com/nomad-space/api
APP?=api
PORT?=80
PORT_APP?=7784
HOST?=api.mvp.nomad.space

RELEASE?=0.0.7
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
CONTAINER_IMAGE?=/nomadspace/${APP}

GOOS?=linux
GOARCH?=amd64

clean:
	rm -f bin/${APP}

dep:
	dep ensure

build: clean dep
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
		-o bin/${APP} \
		-ldflags "-s -w -X ${PROJECT}/src/version.Commit=${COMMIT} \
						-X ${PROJECT}/src/version.Release=${RELEASE} \
						-X ${PROJECT}/src/version.BuildTime=${BUILD_TIME}" \
		./src/cmd/

container: build
	docker build -t $(CONTAINER_IMAGE):$(RELEASE) .

run: container
	docker stop $(CONTAINER_IMAGE):$(RELEASE) || true && docker rm $(CONTAINER_IMAGE):$(RELEASE) || true
	docker run --name ${APP} -p ${PORT}:${PORT} --rm \
		-e "PORT=${PORT}" \
		$(CONTAINER_IMAGE):$(RELEASE)

test:
	go test -v -race ./...

push: container
	docker push $(CONTAINER_IMAGE):$(RELEASE)

minikube: push
	for t in $$(find ./kubernetes -type f -name "*.yaml"); do \
        cat $$t | \
        	sed -E "s/__Release__/$(RELEASE)/g" | \
        	sed -E "s/__ServiceName__/$(APP)/g" | \
        	sed -E "s/__ServiceHost__/$(HOST)/g" | \
        	sed -E "s/__AppPort__/$(PORT_APP)/g" | \
        	sed -E "s/__ContainerImage__/$(CONTAINER_IMAGE2)/g" | \
        	sed -E "s/__ServicePort__/$(PORT)/g"; \
        echo $$"\n"---; \
    done > tmp.yaml;
	kubectl apply -f tmp.yaml
	kubectl patch deployment ${APP} -p "{\"spec\":{\"template\":{\"metadata\":{\"annotations\":{\"date\":\"`date +'%s'`\"}}}}}"