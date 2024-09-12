class fastnode_server () {

  require fastnode_base

  docker::swarm {'cluster_manager':
    init           => true,
  }

  $prefix = "/opt/fastnode-server"

  archive {"${prefix}/fastnode-server.tgz":
    ensure       => present,
    extract      => true,
    source       => "s3://fastnode-deploys/v${$facts['khulnasoft-lab_version']}/fastnode-server.tgz",
    extract_path => '/opt/',
    cleanup      => true,
    creates      => "${prefix}/docker-stack.yml",
    require      => [
      Package['awscli'],
    ],
  }
  -> file {"${prefix}/fastnode-server-deployment-token":
    content => "XXXXXXX\n",
  }
  -> docker::secrets {'fastnode-server-deployment-token':
    secret_name => 'fastnode-server-deployment-token',
    secret_path => "${prefix}/fastnode-server-deployment-token",
  }
  -> docker::stack { 'fastnode-server':
    ensure        => present,
    stack_name    => 'fastnode-server',
    compose_files => ["${prefix}/docker-stack.yml"],
    require       => [Docker::Swarm['cluster_manager']],
  }
}
