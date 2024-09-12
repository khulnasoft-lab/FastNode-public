build: Dockerfile
	docker build -t khulnasoft/build-frontend .

shell: build
	docker run -v $(PWD):/khulnasoft-lab -i -t khulnasoft/build-frontend bash

push: build
	docker push khulnasoft/build-frontend

deploy: push
	ssh build.khulnasoft.com docker pull khulnasoft-lab/build-frontend
