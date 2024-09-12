# == Class: fastnode::postgresql_dev
#
# Sets up postgresql and the community DB (only for development)
#
class fastnode::postgresql_dev {
  class { 'postgresql::globals':
    manage_package_repo => true,
    version             => '9.5',
    encoding            => 'UTF-8',
    locale              => 'en_US.UTF-8',
  }->
  class { 'postgresql::server':
    listen_addresses        => '*',
    pg_hba_conf_defaults    => false,
  }

  postgresql::server::pg_hba_rule { 'local access as postgres user':
    type        => 'local',
    database    => 'all',
    user        => 'postgres',
    auth_method => 'ident',
    order       => 1,
  }

  postgresql::server::pg_hba_rule { 'local access':
    type        => 'local',
    database    => 'all',
    user        => 'all',
    auth_method => 'md5',
    order       => 2,
  }

  postgresql::server::pg_hba_rule { 'IPv4 access':
    type        => 'host',
    database    => 'all',
    user        => 'all',
    address     => '127.0.0.1/32',
    auth_method => 'md5',
    order       => 3,
  }

  postgresql::server::pg_hba_rule { 'IPv6 access':
    type        => 'host',
    database    => 'all',
    user        => 'all',
    address     => '::1/128',
    auth_method => 'md5',
    order       => 3,
  }

  postgresql::server::role { 'fastnode':
    login         => true,
    password_hash => postgresql_password('fastnode', 'fastnode'),
  }

  postgresql::server::db { 'community':
    owner    => 'fastnode',
    user     => 'communityuser',
    password => 'fastnode',
  }

  postgresql::server::db { 'community_test':
    owner    => 'fastnode',
    user     => 'communityuser',
    password => 'fastnode',
  }

  postgresql::server::db { 'apprelease':
    owner    => 'fastnode',
    user     => 'appreleaseuser',
    password => 'fastnode',
  }

  postgresql::server::db { 'localfiles':
    owner    => 'fastnode',
    user     => 'localfilesuser',
    password => 'fastnode',
  }

  postgresql::server::db { 'localfiles_test':
    owner    => 'fastnode',
    user     => 'localfilesuser',
    password => 'fastnode',
  }

  postgresql::server::db { 'localanalysis':
    owner    => 'fastnode',
    user     => 'localanalysisuser',
    password => 'fastnode',
  }

  postgresql::server::db { 'events':
    owner    => 'fastnode',
    user     => 'eventsuser',
    password => 'fastnode',
  }

  postgresql::server::db { 'account':
    owner    => 'fastnode',
    user     => 'accountuser',
    password => 'fastnode',
  }

  postgresql::server::db { 'account_test':
    owner    => 'fastnode',
    user     => 'accountuser',
    password => 'fastnode',
  }
}
