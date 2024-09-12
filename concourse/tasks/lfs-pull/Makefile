.PHONY: all

docker:
	docker build -t khulnasoft/concourse.lfs-pull .

docker-push: docker
	docker push khulnasoft/concourse.lfs-pull
