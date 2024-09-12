import os

from fabric.api import env
from fabric.api import task
from fabric.operations import put, sudo

import docker
import openvpn
import varfastnode
import dns
import puppet
import strongswan

env.user = 'ubuntu'
env.key_filename = '~/.ssh/fastnode-dev.pem'
env.use_ssh_config = True
