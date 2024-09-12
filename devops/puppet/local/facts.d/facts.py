#!/usr/bin/env python3
import os

if 'PUPPET_FACT_fastnode_env' not in os.environ:
    print('fastnode_env=dev')

PREFIX = 'PUPPET_FACT_'

for k, v in os.environ.items():
    if not k.startswith(PREFIX):
        continue
    print('{}={}'.format(k[len(PREFIX):],v))
