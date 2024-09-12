#!/usr/bin/env python
import os
from setuptools import find_packages, setup

def find_scripts(paths):
    for path in paths:
        for parent, dirs, files in os.walk(path):
            if os.path.basename(parent) == 'bin':
                for f in files:
                    yield os.path.join(parent, f)

setup(
    name='fastnode',
    version='0.1.0',
    author='Manhattan Engineering Inc.',
    description='Main python module for Fastnode',
    packages=find_packages(exclude=['tests']),
    scripts=list(find_scripts(["bin", "fastnode"])),

    install_requires=[
        'genshi',  # fastnode.codeexamples
        'sklearn',  # fastnode.ranking
    ],
    extras_require={'test': [
        'pytest',
        'ipdb',
        'numpy',  # tests.ranking
        'numdifftools',  # tests.ranking
    ]}
)
