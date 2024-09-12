# == Class: fastnode::golang::common
#
# Installs some packages that are often useful with golang.
#
class fastnode::golang::common {
  package{ ["mercurial", "protobuf-compiler"]:
    ensure => present,
  }
}
