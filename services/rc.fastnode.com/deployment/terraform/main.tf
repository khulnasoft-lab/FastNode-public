terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
    }
    google = {
      source = "hashicorp/google"
    }
  }
  backend "s3" {
    bucket               = "fastnode-terraform-state"
    workspace_key_prefix = "services/rc.khulnasoft.com"
    key                  = "terraform.tfstate"
    region               = "us-west-1"
  }
}

provider "aws" {
  region = var.aws_region
}

provider "google" {
  region  = var.gcp_region
  project = var.gcp_project
}

provider "google-beta" {
  region  = var.gcp_region
  project = var.gcp_project
}

resource "google_compute_global_address" "svc" {
  name = "rc-fastnode-com-${terraform.workspace}"
}

resource "google_service_account" "default" {
  account_id   = "svc-rc-fastnode-com-${terraform.workspace}"
  display_name = "svc-rc-fastnode-com-${terraform.workspace}"
}

resource "google_project_iam_member" "default" {
  project = var.gcp_project
  role    = "roles/iam.serviceAccountTokenCreator"
  member  = "serviceAccount:svc-rc-fastnode-com-${terraform.workspace}@${var.gcp_project}.iam.gserviceaccount.com"
}

resource "google_service_account_iam_binding" "default" {
  service_account_id = google_service_account.default.name
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "serviceAccount:${var.gcp_project}.svc.id.goog[rc-fastnode-com-${terraform.workspace}/service]",
  ]
}

resource "aws_iam_role" "role" {
  name = "svc-rc-fastnode-com-${terraform.workspace}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "accounts.google.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "accounts.google.com:sub": "${google_service_account.default.unique_id}"
        }
      }
    }
  ]
}
EOF
}


data "aws_iam_policy_document" "policy" {
  statement {
    sid       = "1"
    actions   = ["s3:ListAllMyBuckets"]
    resources = ["arn:aws:s3:::*"]
  }

  statement {
    sid       = "2"
    actions   = ["s3:ListBucket"]
    resources = ["arn:aws:s3:::fastnode-metrics/*"]
  }

  statement {
    sid       = "3"
    actions   = ["s3:GetObject"]
    resources = ["arn:aws:s3:::fastnode-metrics/enrichment/maxmind/raw/country/latest/*"]
  }

  statement {
    sid       = "4"
    actions   = ["s3:GetBucketLocation"]
    resources = ["arn:aws:s3:::fastnode-metrics"]
  }
}

resource "aws_iam_policy" "policy" {
  name   = "svc-rc-fastnode-com-${terraform.workspace}"
  path   = "/"
  policy = data.aws_iam_policy_document.policy.json
}

resource "aws_iam_role_policy_attachment" "default_attachment" {
  role       = aws_iam_role.role.name
  policy_arn = aws_iam_policy.policy.arn
}
