resource "aws_key_pair" "fastnode-prod" {
  key_name   = "fastnode-prod"
  public_key = "ssh-rsa XXXXXXX"
}
