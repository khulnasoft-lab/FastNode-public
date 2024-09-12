variable service_name {
  default = "t-fastnode-com"
}

variable region {
  default = "us-east-1"
}

variable webserver_repository_name {
  type = string
  default = "t-fastnode-com-webserver"
}

variable fluentd_repository_name {
  type = string
  default = "t-fastnode-com-fluentd"
}

variable tag {
  type = string
}

variable webserver_port {
  default = 9000
}

variable cpu {
  default = 1
}

variable memory {
  default = 2
}
