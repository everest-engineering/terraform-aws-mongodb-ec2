provider "aws" {
  region  = var.region
  profile = "terraform-provisioner-ansible"
}

module "mongodb" {
  source          = "../../"
  vpc_id          = var.vpc_id
  subnet_id       = var.subnet_id
  instance_type   = "t2.micro"
  ssh_user        = "ubuntu"
  ami_filter_name = "ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"
  ami_owners      = ["099720109477"]
  mongodb_version = "4.2"
  replicaset_name = "mongo-rp0"
  replica_count   = 1
  data_volumes    = var.data_volumes
  private_key     = file("~/.ssh/id_rsa")
  public_key      = file("~/.ssh/id_rsa.pub")
  bastion_host    = var.bastion_host
  tags = {
    Name        = "MongoDB Server"
    Environment = "terraform-mongo-testing"
  }
}

variable "region" {
  type        = string
  description = "AWS Region"
}

variable "vpc_id" {
  type        = string
  description = "VPC Id"
}

variable "subnet_id" {
  type        = string
  description = "Subnet Id"
}

variable "data_volumes" {
  type = list(object({
    ebs_volume_id     = string
    availability_zone = string
  }))
  description = "List of EBS volumes"
}

variable "bastion_host" {
  type        = string
  description = "Bastion host Public IP"
}

output "mongo_server_ip_address" {
  value = module.mongodb.mongo_server_private_ip
}
