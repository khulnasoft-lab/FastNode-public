class release_server {
  require fastnode_base

  case $facts['virtual'] {
    'docker': {
      include release_server::service_initd
    }
    default:  {
      include release_server::service_systemd
    }
  }

  file {'/var/fastnode/bin/run.sh':
    content => epp('release_server/run.sh.epp', {
      executable  => '/var/fastnode/bin/release server',
      env_config  => {
        'RELEASE_DB_DRIVER' => 'postgres',
        'ROLLBAR_ENV'       => 'production',
      },
      secret_keys => ['RELEASE_DB_URI', 'ROLLBAR_TOKEN']
    }),
    owner   => fastnode,
    group   => fastnode,
    mode    => '0755',
    notify  => Service['release_server'],
    require => File['/var/fastnode/bin'],
  }
  -> archive { '/var/fastnode/bin/release':
    ensure => present,
    source => "s3://fastnode-deploys/v${$facts['khulnasoft-lab_version']}/release"
  } ~> exec {'khulnasoft-lab permissions':
    command     => 'chmod 755 /var/fastnode/bin/release',
    path        => ['/bin', '/usr/bin'],
    refreshonly => true,
    onlyif      => "test `stat -c '%a' /var/fastnode/bin/release` != 755",
    notify      => Service['release_server'],
  }
}
