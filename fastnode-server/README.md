# Fastnode Team Server

## Setup
Currently, only Ubuntu 20.04 and RHEL 7 variants are officially supported.

The `nvidia-docker` image in GCP can be used to bootstrap steps 1-3 below.

1.  Run `make fastnode-server.tgz` from this directory.
2.  Transfer the archive to the host and extract its contents, which should result in a `fastnode-server` directory.
3.  Run the `configure` script for the host OS. The script should install Docker and all dependencies.
4.  Optional: Perform the [Postinstall Steps](https://docs.docker.com/engine/install/linux-postinstall/) so Docker commands can be run without `sudo`
5.  Create a `deployment_token` with `openssl rand -hex 32 | sudo docker secret create fastnode-server-deployment-token -`
6.  Run `docker swarm init`
7.  Run `docker stack deploy -c docker-stack.yml fastnode-server`
8.  To check that all services have spun up, check the output of `docker service ls`:
    ```
    ID                  NAME                             MODE                REPLICAS            IMAGE                                          PORTS
    XXXXXXX        fastnode-server_edge-ingress         replicated          1/1                 khulnasoft-lab/fastnode-server-edge-ingress:latest         *:8500->8500/tcp, *:9902->9901/tcp
    XXXXXXX        fastnode-server_metadata             replicated          1/1                 khulnasoft-lab/fastnode-server-metadata:latest             *:8080->8080/tcp
    XXXXXXX        fastnode-server_models               replicated          1/1                 khulnasoft-lab/fastnode-server-models:latest
    XXXXXXX        fastnode-server_models-ingress       replicated          1/1                 khulnasoft-lab/fastnode-server-models-ingress:latest       *:9901->9901/tcp
    XXXXXXX        fastnode-server_models-stats-proxy   replicated          1/1                 khulnasoft-lab/fastnode-server-models-stats-proxy:latest   *:8601->8601/tcp
    ```

## Known Tokens

Known tokens can be found in the [KTS Deployment Token Google Sheet](https://docs.google.com/XXXXXXX)

### Deployment IDs

To derive a deployment ID from a token, use the following Python 3 snippet:

```python3
import hashlib
import base64

token = "XXXXXXX"
print(base64.b64encode(hashlib.sha256((token+'\n').encode()).digest()))
```
