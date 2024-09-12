BUILD_DIR := /tmp/import-exploration-image

SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

image:
	# Could also use setup.py sdist but we want to avoid getting the host python interpreter involved
	rm -f fastnode-python.tar
	tar cf $(SELF_DIR)/fastnode-python.tar \
		--exclude build \
		--exclude lib \
		--exclude dist \
		--exclude __pycache__ \
		--exclude '*.pyc' \
		--exclude fastnode-python.egg \
		-C $(SELF_DIR)/../.. \
		fastnode-python
	docker build -t khulnasoft/import-exploration $(SELF_DIR)

push: build
	docker push khulnasoft/import-exploration
