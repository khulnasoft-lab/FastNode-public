import os

_BOOTSTRAP_TEMPLATE = """#!/bin/bash
aws s3 cp $S3ROOT/bootstrap/requirements-emr.txt ./
aws s3 cp $S3ROOT/bootstrap/fastnode.tar.gz ./
sudo alternatives --set python /usr/bin/python2.7
sudo yum-config-manager --enable epel
sudo pip install -r requirements-emr.txt
sudo pip install fastnode.tar.gz

sudo mkdir -p /var/fastnode
sudo chown -R hadoop:hadoop /var/fastnode
mkdir -p /var/fastnode/s3cache/tmp
"""


def template_with_root(root):
    return (_BOOTSTRAP_TEMPLATE.replace("$S3ROOT", root))
