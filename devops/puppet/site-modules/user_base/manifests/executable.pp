class user_base::executable {
  file {'/mnt/fastnode':
      ensure => directory,
      owner  => ubuntu,
      group  => ubuntu
    }

    $process_file = "${facts[khulnasoft-lab_version]}/${user_base::process_name}"

    file {'/mnt/fastnode/releases':
      ensure => directory,
      owner  => ubuntu,
      group  => ubuntu
    }

    archive { "/mnt/fastnode/releases/${process_file}":
      ensure => present,
      source => "s3://fastnode-deploys/${process_file}"
    } ~> exec {'khulnasoft-lab permissions':
      command     => "chmod 755 /mnt/fastnode/releases/${process_file}",
      path        => ['/bin', '/usr/bin'],
      refreshonly => true,
      onlyif      => "test `stat -c '%a' /mnt/fastnode/releases/${process_file}` != 755"
    }

    file {'/mnt/fastnode/s3cache':
      ensure => directory,
      owner  => ubuntu,
      group  => ubuntu
    }

    file {'/mnt/fastnode/logs':
      ensure => directory,
      owner  => ubuntu,
      group  => ubuntu
    }

    file {'/mnt/fastnode/certs':
      ensure => directory,
      owner  => ubuntu,
      group  => ubuntu
    } -> archive { '/mnt/fastnode/certs/rds-combined-ca-bundle.pem':
      ensure => present,
      source => 's3://XXXXXXX/rds-combined-ca-bundle.pem'
    }

    file {'/mnt/fastnode/tmp':
      ensure => directory,
      owner  => ubuntu,
      group  => ubuntu
    }

    archive { '/var/fastnode/config.sh':
      ensure => present,
      source => "s3://XXXXXXX/config/${facts['aws_region']}.sh"
    } ~> exec {'khulnasoft-lab_config_permissions':
      command     => 'chmod 755 /var/fastnode/config.sh',
      path        => ['/bin', '/usr/bin'],
      refreshonly => true,
      onlyif      => "test `stat -c '%a' /var/fastnode/config.sh` != 755"
    }

    file {'/var/fastnode':
      ensure => link,
      target => '/mnt/fastnode',
      owner  => ubuntu,
      group  => ubuntu
    }

    archive { '/usr/local/libtensorflow-cpu-linux-x86_64-1.15.0.tar.gz':
      ensure       => present,
      source       => 's3://fastnode-data/tensorflow/libtensorflow-cpu-linux-x86_64-1.15.0.tar.gz',
      extract      => true,
      extract_path => '/usr/local',
    } ~> exec {'ldconfig':
      command     => 'ldconfig',
      path        => ['/sbin'],
      refreshonly => true,
      notify      => Service[$user_base::process_name],
    }

    file {'/var/fastnode/run.sh':
      content => epp('user_base/run.sh.epp', {executable=>"/var/fastnode/releases/${process_file}"}),
      owner   => ubuntu,
      group   => ubuntu,
      mode    => '0755',
    } ~> systemd::unit_file { "${user_base::process_name}.service":
      content => file('user_base/khulnasoft-lab.service'),
    }
    ~> service {$user_base::process_name:
      ensure   => 'running',
      provider => 'systemd',
    }
}
