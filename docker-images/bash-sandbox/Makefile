build: Dockerfile
	docker build -t khulnasoft/bashsandbox .

push: build
	docker push khulnasoft/bashsandbox

deploy: push
	ssh curation.khulnasoft.com docker pull khulnasoft-lab/bashsandbox
