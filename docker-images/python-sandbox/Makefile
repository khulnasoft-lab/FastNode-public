build: Dockerfile
	docker build -t khulnasoft/pythonsandbox .

push: build
	docker push khulnasoft/pythonsandbox

deploy: push
	ssh curation.khulnasoft.com docker pull khulnasoft-lab/pythonsandbox
