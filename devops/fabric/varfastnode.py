from datetime import datetime

from fabric.api import run
from fabric.api import task
from fabric.contrib import files
from fabric.operations import get, put


@task
def backup():
    ts = datetime.now().strftime("%Y-%m-%d-%H-%M-%S")
    run("tar -cvzf varfastnode-backup-%s.tar.gz /var/fastnode" % ts)
    get(remote_path="varfastnode-backup-%s.tar.gz" % ts, local_path="backup/")
    run("rm varfastnode-backup-*.tar.gz")
