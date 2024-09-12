DEV_TAG=airflow_monetizable_dev:latest
KHULNASOFT-LAB=$(PWD)/../../..
DOCKER_RUN_CMD=docker run --rm -it -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) -v $(KHULNASOFT-LAB):/go/src/github.com/khulnasoft-lab/fastnode $(DEV_TAG)

ECR_REPO_URL=XXXXXXX.dkr.ecr.us-west-1.amazonaws.com
ECR_PACKAGE_NAME=fastnode-airflow-monetizable
TAG=$(shell git rev-parse --short HEAD)

KHULNASOFT-LAB=$${PWD%/khulnasoft-lab/**}/khulnasoft-lab
CWD_RELATIVE=$${PWD\#/**/khulnasoft-lab}
GO_IMAGE=golang:1.15.3-buster

dev.build: build/monetizable
	docker build --build-arg BUILD_HASH=$(TAG) -t $(DEV_TAG) .

dev.shell:
	@exec $(DOCKER_RUN_CMD) /bin/bash

docker.login:
	aws ecr get-login-password --region us-west-1 | docker login --username AWS --password-stdin $(ECR_REPO_URL)

docker.build: build/monetizable
	docker build --build-arg BUILD_HASH=$(TAG) -t $(ECR_PACKAGE_NAME):$(TAG) .

docker.tag:
	docker tag $(ECR_PACKAGE_NAME):$(TAG) $(ECR_REPO_URL)/$(ECR_PACKAGE_NAME):$(TAG)

docker.push:
	docker push $(ECR_REPO_URL)/$(ECR_PACKAGE_NAME):$(TAG)

docker.all: docker.login docker.build docker.tag docker.push

build/monetizable: build main.go
	docker run --rm -e "GOPRIVATE=github.com/khulnasoft-lab/*" \
		-v $(KHULNASOFT-LAB):/go/src/github.com/khulnasoft-lab/fastnode \
		-v $(PWD)/build:/build \
		-w /go/src/github.com/khulnasoft-lab/$(CWD_RELATIVE) \
		$(GO_IMAGE) go build -o /build/monetizable .

build:
	mkdir -p build
