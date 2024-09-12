#!/usr/bin/env python
import sys
from setuptools import find_packages, setup

requirements = []


setup(
    name='fastnode.ml',
    version='0.1.0',
    author='Manhattan Engineering Inc.',
    description='Fastnode Python ML',
    packages=find_packages(exclude=['tests']),
    install_requires=requirements,
)
