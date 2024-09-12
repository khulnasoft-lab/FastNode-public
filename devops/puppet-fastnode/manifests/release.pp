# == Class: fastnode::release
#
# Configuration for release.khulnasoft.com, Fastnode's release server, as well as
# the staging and mock release servers.

class fastnode::release (
  $environment = undef,
  $hostname = undef,
) {
  if str2bool($::vagrant) {
    $owner = "vagrant"
    $group = "vagrant"
  } else {
    $owner = "ubuntu"
    $group = "ubuntu"
  }
  
  file { ["/var/fastnode", "/var/fastnode/releases", "/var/fastnode/log"]:
    ensure => directory,
    owner => $owner,
    group => $group,
  }

  include fastnode::python # for pip
  include nginx
  include fastnode::ubuntu::bootstrap

  # Set the system environment variables
  file { "/etc/environment":
    content => template("fastnode/prod/environment.sh.erb"),
    owner   => "root",
    group   => "root",
  }

  # install s3cmd (for deployments)
  package { 's3cmd':
    provider => 'pipx',
    ensure => present,
  }

  if str2bool($::vagrant) {
    exec { "make-certs":
      command => "/usr/bin/openssl req -new -newkey rsa:2048 -days 364 -nodes -x509 -subj '/C=US/ST=CA/L=SF/O=./OU=./CN=192.168.30.10' -keyout /etc/ssl/release.khulnasoft.com.key -out /etc/ssl/server.crt",
      unless  => "/bin/ls /etc/ssl/server.crt",
    }
  } else {
    exec { "get-cert":
      command => "/usr/local/bin/s3cmd get s3://XXXXXXX/ssl/server.crt /etc/ssl/server.crt",
      unless => "/bin/ls /etc/ssl/server.crt",
    }
    exec { "get-key":
      command => "/usr/local/bin/s3cmd get s3://XXXXXXX/ssl/server.key /etc/ssl/server.key",
      unless => "/bin/ls /etc/ssl/server.key",
    }
  }

  # nginx config
  file { "/etc/nginx/sites-available/release.khulnasoft.com":
    content => template("fastnode/nginx/release.khulnasoft.com.erb"),
    owner => "root",
    group => "root",
    notify => Service["nginx"],
  } ->
  file { "/etc/nginx/sites-enabled/release.khulnasoft.com":
    ensure => 'link',
    target => "/etc/nginx/sites-available/release.khulnasoft.com",
    notify => Service["nginx"],
  }

  package { "ca-certificates":
    ensure => present,
  }
}
