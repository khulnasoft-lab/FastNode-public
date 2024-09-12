resource "aws_s3_bucket" "localcontent-bucket" {
  bucket = "fastnode-local-content-${var.region}"
}

resource "aws_s3_bucket" "localsymbols-bucket" {
  bucket = "fastnode-local-symbols-${var.region}"
}

resource "aws_s3_bucket" "prod-data-bucket" {
  bucket = "fastnode-prod-data-${var.region}"
}
