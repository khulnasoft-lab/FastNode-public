#!/usr/bin/env python
import sys
from setuptools import find_packages, setup

setup(
    name='fastnode.metrics',
    version='0.1.0',
    author='Manhattan Engineering Inc.',
    description='Fastnode Metrics',
    packages=find_packages(),
    install_requires=[
        "jinja2>=2",
        "PyYAML>=5",
        "click>=7",
    ],
    entry_points = {
        'console_scripts': ['fastnode-metrics-schemas=fastnode_metrics.json_schema:main'],
    },
    python_requires='>=3.6',
    include_package_data = True,
)
