# == Class: fastnode::curation
#
# Base module for curation configuration.
#
class fastnode::curation {
  # Configuration of nginx happens via hiera
  include nginx
  include fastnode::python

  # Installs dependencies for Stanford's CoreNLP parser
  include fastnode::curation::corenlp

  # Required by the codeexample authoring tool's linter
  package { 'autopep8':
    provider => 'pipx',
    ensure   => present,
  }
  package { 'pylint':
    provider => 'pipx',
    ensure   => present,
  }
}
