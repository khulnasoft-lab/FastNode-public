variable "gcp_regions" {
  default = ["us-west1"]
}

variable "gcp_project" {
  default = "fastnode-prod-XXXXXXX"
}

variable aws_region {
  default = "us-west-1"
}

variable "versions" {
  type = map(string)
}
