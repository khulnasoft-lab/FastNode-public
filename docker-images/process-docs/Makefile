OUTPUT := process-docs.tar

image: Dockerfile
	docker build -t khulnasoft/process-docs .

save: image
	docker save khulnasoft-lab/process-docs > $(OUTPUT) 
