node "fastnode-base.fastnode.dev" {
  include fastnode_base
}

node "metrics-collector.fastnode.dev" {
  include fastnode_base
  include metrics_collector
}

node "airflow.fastnode.dev" {
  include airflow
}


node "user-node.fastnode.dev" {
  include user_base
}

node "user-mux.fastnode.dev" {
  include user_base
}

node /release-[0-9a-z]+/ {
  include release_server
}

node /nchan-[0-9a-z]+/ {
  include nchanted
}

node /fastnodeserver-[0-9a-z]+/ {
  include fastnode_server
}
