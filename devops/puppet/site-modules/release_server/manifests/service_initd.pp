class release_server::service_initd {
  initscript::service { 'release_server':
    cmd            => '/var/fastnode/bin/run.sh',
    define_service => true,
  }
}
