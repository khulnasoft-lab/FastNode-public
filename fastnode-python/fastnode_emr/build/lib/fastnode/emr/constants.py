import os

PIPELINE_FILE = 'pipeline.yaml'
# FASTNODE_COMMON_ROOT is the path to fastnode-python/fastnode_emr in the fastnodeco repo
FASTNODE_PYTHON = os.environ.get("FASTNODE_EMR_ROOT",
                             os.path.join(os.environ["GOPATH"], "src/github.com/khulnasoft-lab/fastnode/fastnode-python/fastnode_emr"))
FASTNODE_EMR_BUCKET = 'fastnode-emr'
BUNDLE_DIR = 'bundle'
