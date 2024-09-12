# == Class: fastnode::curation::corenlp
#
# Installs some packages that are required for Standford's CoreNLP parser
#
class fastnode::curation::corenlp {
  include fastnode::java

  # Some python libraries required for the python interface to corenlp.
  package { ['pexpect', 'unidecode', 'jsonrpclib']:
    provider => 'pipx',
    ensure   => present,
  }
}
