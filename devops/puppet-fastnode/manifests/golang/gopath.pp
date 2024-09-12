# == Class: fastnode::golang::gopath
#
# Sets up a user's $GOPATH, including src, bin, pkg and the environment variable.
#
# === Parameters:
#
# [*path*]
#   The path to use for $GOPATH
#
# [*owner*]
#   The owner of $GOPATH
#
# [*group*]
#   Group of $GOPATH
#
class fastnode::golang::gopath(
  $path   = "/usr/local/gopath",
  $owner  = undef,
  $group  = undef,
) {
  include fastnode::golang

  file { $path:
    ensure => directory,
    owner  => $owner,
    group  => $group,
  }

  file { [ "$path/src", "$path/bin", "$path/pkg" ]:
    ensure => directory,
    owner  => $owner,
    group  => $group,
  }

  file { "/etc/profile.d/gopath.sh":
    content => template("fastnode/golang/gopath.sh.erb"),
    owner   => root,
    group   => root,
    mode    => "a+x",
  }

}
