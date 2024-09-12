.PHONY: deployment-token clean

deployment-token:
	openssl rand -hex 32

deployment-id:
	python -c 'import hashlib; import base64; print(base64.b64encode(hashlib.sha256(("${TOKEN}\n").encode()).digest()))'

clean:
	rm -rf fastnode-teams-ctl fastnode-server*

fastnode-teams-ctl:
	GOOS=linux go build ./tunable-models/cmds/fastnode-team-ctl


tag-all:
ifeq ($(TAG),)
	$(error Must set TAG to push)
endif
	cd edge-ingress && make TAG=$(TAG) build push
	cd metadata && make TAG=$(TAG) clean build push
	cd models-ingress && make TAG=$(TAG) build push
	cd models-stats-proxy && make TAG=$(TAG) clean build push
	cd tunable-models && make TAG=$(TAG) build push

VERSION = v0.6.0
DIR = fastnode-server-$(VERSION)

fastnode-server.tgz: clean fastnode-teams-ctl
	# setup deployment token
ifeq ($(TOKEN),)
	$(error Must set TOKEN to build fastnode-server.tgz)
endif

	mkdir $(DIR)
	# to store tuned models, mounted by tunable-models
	mkdir $(DIR)/tuned-models
	# to store code to tune on, mounted by tunable-models
	mkdir $(DIR)/repositories
	# to store fastnode-team-ctl binary
	mkdir $(DIR)/bin

	echo $(TOKEN) > $(DIR)/deployment_token
	@echo building with deployment_token $(TOKEN)

	cp fastnode-team-ctl $(DIR)/bin/
	cp nvidia-docker/configure_*.sh $(DIR)
	cp nvidia-docker/launch_kts.sh $(DIR)
	cp docker-stack.yml $(DIR)
	tar cvzf $@ $(DIR)
	rm -r $(DIR) fastnode-team-ctl
