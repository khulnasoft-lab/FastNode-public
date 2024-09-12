# == Class: fastnode::docker
#
# Sets up docker
#
class fastnode::docker (
  $user = [],
) {
  package { 'docker-engine':
    ensure  => present,
  } ->
  exec { "docker-system-user-${user}":
    command  => "/usr/sbin/usermod -aG docker ${user}",
    unless =>  "/bin/cat /etc/group | grep '^docker:' | grep -qw ${user}",
  }
}
