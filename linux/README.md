# Developing on Linux

Below are a set of instructions that get you from a vanilla Linux (Ubuntu) installation to building Fastnode.

You will have to generate SSH keys, setup GOPATH and clone the repository from the `Get the codebase` section below.

# Setting up to develop on Linux

This document will go through how to setup a vanilla Linux (Ubuntu) installation to work with the Fastnode codebase. Once you complete this guide, you should be able to build both the Copilot and fastnoded.exe, and run them together by running the `./run_fastnoded.sh` command in this directory.

## Install Golang
Go to the [official golang website](https://golang.org/dl/) and download the latest version of Golang supported by our codebase.
You can also use [`devops/scripts/install-golang.sh`](https://github.com/khulnasoft-lab/fastnode/blob/master/devops/scripts/install-golang.sh)

## Get the codebase

Generate your ssh keys:

```sh
ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
```

Add your **PUBLIC KEY** (e.g `cat ~/.ssh/id_rsa.pub`) it to your GitHub account under Settings > SSH Keys.

Setup your `GOPATH` and clone the repository
```
mkdir -p ~/go/src/github.com/khulnasoft-lab
cd ~/go/src/github.com/khulnasoft-lab
git clone git@github.com:khulnasoft-lab/khulnasoft-lab
```

Then, run
```
git lfs pull
```
to grab the necessary binary files.

## Start building

Go to the `linux` directory.

For 1st time setup, you'll need to build the copilot first:

```sh
./build_electron.sh
```

To run Fastnode:

```sh
./run_fastnoded.sh
```

If the visual copilot doesn't startup automatically, run:

```sh
./linux-unpacked/fastnode
```

If you get this error:
```sh
go build github.com/khulnasoft-lab/fastnode/vendor/github.com/kiteco/tensorflow/tensorflow/go: invalid flag in #cgo LDFLAGS: -Wl,-rpath,/home/fastnode/go/src/github.com/khulnasoft-lab/fastnode/vendor/github.com/kiteco/tensorflow/tensorflow/go/../../../../../../linux/tensorflow/lib,-z,undefs
```
Add `export CGO_LDFLAGS_ALLOW=".*"` to your `.bashrc` or `.zshrc` file.

## VM Notes

If you're running Linux within a VM, and would like to mount a network drive to share code, the default directory is `/mnt/hgfs/`.
