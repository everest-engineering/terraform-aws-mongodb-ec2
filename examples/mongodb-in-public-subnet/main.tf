provider "aws" {
  region  = "us-east-1"
  profile = "terraform-provisioner-ansible"
}

data "aws_vpc" "default" {
  default = true
}

data "aws_subnet" "subnet" {
  vpc_id            = data.aws_vpc.default.id
  availability_zone = var.availability_zone
}

module "mongodb" {
  source            = "..\/..\/ec2-replicaset"
  vpc_id            = data.aws_vpc.default.id
  subnet_id         = data.aws_subnet.subnet.id
  ssh_user          = "ubuntu"
  instance_type     = "t2.micro"
  ami_filter_name   = "ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"
  ami_owners        = ["099720109477"]
  data_volumes      = var.data_volumes
  mongodb_version   = "4.2"
  replicaset_name   = "mongo-rp0"
  replica_count     = 1
  private_key       = file("~/.ssh/id_rsa")
  public_key        = file("~/.ssh/id_rsa.pub")
  tags = {
    Name        = "MongoDB Server"
    Environment = "terraform-mong-testing"
  }

}

variable "data_volumes" {
  type = list(object({
    ebs_volume_id = string
    availability_zone = string
  }))
  description = "List of EBS volumes"
  default = []
}
variable "availability_zone" {
  type    = string
  default = "us-east-1a"
}

output "mongo_server_ip_address" {
  value = module.mongodb.mongo_server_public_ip
}
