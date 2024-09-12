image:
	GOOS=linux go build github.com/khulnasoft-lab/fastnode/fastnode-go/cmds/local-code-worker
	docker build -t khulnasoft/local-code-worker .

push:
	docker push khulnasoft/local-code-worker

pull:
	docker pull khulnasoft-lab/local-code-worker
	docker pull khulnasoft-lab/process-docs

run:
	mkdir -p $(LOCAL_WORKER_STORAGE)
	touch $(SLOT_DB_URI)
	docker run -it \
		-e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) \
		-e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) \
		-e SLOT_DB_URI=/slots.db \
		-e SLOT_DB_DRIVER=$(SLOT_DB_DRIVER) \
		-e LOCAL_WORKER_STORAGE=/files \
		-e LOCAL_WORKER_REQUEST_QUEUE=$(LOCAL_WORKER_REQUEST_QUEUE) \
		-e LOCAL_WORKER_RESPONSE_QUEUE=$(LOCAL_WORKER_RESPONSE_QUEUE) \
		-e LOCAL_WORKER_BUCKET=$(LOCAL_WORKER_BUCKET) \
		-v $(SLOT_DB_URI):/slots.db \
		-v $(LOCAL_WORKER_STORAGE):/files \
		khulnasoft-lab/local-code-worker
