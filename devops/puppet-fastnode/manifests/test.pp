# == Class: fastnode::test
#
# Common configuration for test machines
#
class fastnode::test (
  $environment = undef,
  $hostname = undef,
  $vagrant_ip = undef,
) {
  # If we are in vagrant the owner and group is "vagrant"
  # Otherwise, the owner and group is "ubuntu"
  if str2bool($::vagrant) {
    $owner = "vagrant"
    $group = "vagrant"
  } else {
    $owner = "ubuntu"
    $group = "ubuntu"
  }

  # Install things
  include nginx
  include fastnode::ubuntu::bootstrap
  include fastnode::golang::install
  include fastnode::postgresql_dev

  # Setup /var/fastnode to point to /mnt/fastnode. In vagrant, this is
  # not really needed, but it makes it consistent with our EC2 nodes
  # for now (which mount their storage at /mnt)
  file { "/mnt/fastnode":
    ensure => directory,
    owner  => $owner,
    group  => $group,
  } ->
  file { "/var/fastnode":
    ensure => link,
    owner  => $owner,
    group  => $group,
    target => "/mnt/fastnode"
  } ->
  file { ["/var/fastnode/log", "/var/fastnode/data", "/var/fastnode/bin", "/var/fastnode/tmp"]:
    ensure => directory,
    owner  => $owner,
    group  => $group
  }

  file { "/deploy":
    ensure => directory,
    owner  => $owner,
    group  => $group
  }

  # Create self-signed ssl certificates
  file { "$home/certs":
    ensure => directory,
    owner  => $owner,
    group  => $group,
  }

  if str2bool($::vagrant) {
    exec { "make-certs":
      command => "/usr/bin/openssl req -new -newkey rsa:2048 -days 364 -nodes -x509 -subj '/C=US/ST=CA/L=SF/O=./OU=./CN=$vagrant_ip' -keyout $home/certs/server.key -out $home/certs/server.crt",
      unless  => "/bin/ls $home/certs/server.key",
      require => File["$home/certs"],
    }
  } else {
    exec { "make-certs":
      command => "/usr/bin/yes 'xx' | /usr/bin/openssl req -new -newkey rsa:2048 -days 364 -nodes -x509 -keyout $home/certs/server.key -out $home/certs/server.crt",
      unless  => "/bin/ls $home/certs/server.key",
      require => File["$home/certs"],
    }
  }

  # Set the system environment variables
  file { "/etc/environment":
    content => template("fastnode/prod/environment.sh.erb"),
    owner   => "root",
    group   => "root",
  }

  # Set nginx configuration
  file { "/etc/nginx/sites-available/usernode":
    content => template("fastnode/nginx/usernode.erb"),
    owner   => "root",
    group   => "root",
    notify  => Service["nginx"],
  } ->
  file { "/etc/nginx/sites-enabled/usernode":
    ensure => 'link',
    target => "/etc/nginx/sites-available/usernode",
    notify  => Service["nginx"],
  }
}
