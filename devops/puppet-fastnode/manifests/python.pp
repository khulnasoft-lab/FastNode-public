# == Class: fastnode::python
#
# Sets up python, pip, and virualenv. Currently via counsyl/python.
#
class fastnode::python {
  include '::python'
  include '::python::devel'
  include '::python::virtualenv'

  # dev headers for python3.4
  package { 'python3.4-dev':
    ensure => present,
  }

  # requirements for various packages installed by fastnode-python module
  package { ['pkg-config', 'gfortran', 'libopenblas-dev', 'liblapack-dev', 'libfreetype6-dev', 'libxft-dev']:
    ensure => present,
  }
}
