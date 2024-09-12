Cloud Development
=================

In the past, Fastnode's language analysis facilities ran on an AWS/Azure backend instead of the users machine.
There are still several backend components including exposing the symbol API for serving web docs as well as
servicing the web sandbox.


### Vagrant

We use VMs for backend development to guarantee a consistent environment between development and production.
To get this set up, first [set up Vagrant](vagrant-boxes/fastnode-dev/README.md)

Once you have a shell in the virtual machine, the khulnasoft-lab repo's working directory should be at:

```sh
$HOME/go/src/github.com/khulnasoft-lab/fastnode
```

NOTE: This is a symlink to `/khulnasoft-lab`, mounted as a NFS share in the `Vagrantfile`

All commands (`make *`, `go build`, etc) must be run from the full `$HOME/go/src/github.com/khulnasoft-lab/fastnode` path (not a symlinked directory).
From here, you may need to repeat some of the steps from the original dev setup, e.g:

```sh
# Install dependencies
make install-deps
```

Because `user-node` takes too many resources to load locally, there are test instances available on AWS/Azure for you to run/test your development changes to `user-node`.
Please see https://fastnode.quip.com/Phk4AB8lLqh9 for a list of test instances; we no longer have per-developer test instances,
so please notify others before deploying the backend or otherwise running resource intensive processes.

Once you have your test instance, you can deploy your local changes to it by running:

```sh
cd ~/go/src/github.com/khulnasoft-lab/fastnode
./scripts/deploy_test.sh test-N.khulnasoft.com
```


#### Infrastructure

Our AWS infrastructure makes use of Terraform (http://www.terraform.io (http://www.terraform.io/)). Terraform helps us manage our AWS topology. Please do not modify this unless you know what you are doing :). Our terraform configuration files can be found in github.com/khulnasoft-lab/fastnode/devops/terraform (http://github.com/khulnasoft-lab/fastnode/devops/terraform).

We use Fabric to execute some commands on remote hosts (others are simply shell scripts that invoke SSH). The fabric scripts can be found at github.com/khulnasoft-lab/fastnode/devops/fabric (http://github.com/khulnasoft-lab/fastnode/devops/fabric).
