class fastnode_base {
  $puppet_root = "/opt/fastnode/puppet"

  class { '::logrotate':
    ensure => 'latest',
    config => {
      dateext       => true,
      compress      => true,
      delaycompress => true,
      minsize       => "10M",
      rotate        => 5,
      ifempty       => true,
    }
  }

  file {['/var/fastnode', '/var/fastnode/aws', '/var/fastnode/bin']:
    ensure => directory,
  }

  if $gce != undef {
    file {'/var/fastnode/aws/credentials':
      content => epp('fastnode_base/var/fastnode/aws/gcp_credentials'),
      mode    => '0777',
      require => File['/var/fastnode/aws']
    }
    file {'/var/fastnode/aws/config':
      content => file('fastnode_base/var/fastnode/aws/gcp_config'),
    }

    file {'/root/.aws':
      ensure => directory
    }
    ~> file {'/root/.aws/config':
      ensure  => link,
      target  => '/var/fastnode/aws/config',
      require => File['/var/fastnode/aws/config']
    }
  }

  file {'/var/fastnode/aws/run_with_secrets':
    content => file('fastnode_base/var/fastnode/aws/run_with_secrets'),
    mode    => '0777',
    require => File['/var/fastnode/aws']
  }

  file {'/etc/modprobe.d':
    ensure => directory
  }

  file {'/etc/modprobe.d/blacklist-conntrack.conf':
    content => 'blacklist nf_conntrack'
  }

  file {'/etc/sysctl.d':
    ensure => directory
  }

  file {'/etc/sysctl.d/60-fastnode.conf':
    content => file('fastnode_base/etc/sysctl.d/60-fastnode.conf')
  }

  file {'/usr/local/bin/puppet':
    ensure => 'link',
    target => '/opt/puppetlabs/bin/puppet'
  }

  file {'/usr/local/bin/facter':
    ensure => 'link',
    target => '/opt/puppetlabs/bin/facter'
  }

  include 'fastnode_base::packages'
  include 'fastnode_base::users'
  include 'fastnode_base::monitoring'
}
