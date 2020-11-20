variable "vpc_id" {
  type        = string
  description = "VPC Id"
}

variable "subnet_id" {
  type        = string
  description = "Subnet Id"
}

variable "instance_type" {
  type        = string
  description = "AWS EC2 instance type"
  default     = "t2.micro"
}

variable "ami_id" {
  type        = string
  description = "AWS AMI Id"
  default     = ""
}

variable "ami_filter_name" {
  type        = string
  description = "AWS AMI Name filter value"
  default     = "ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"
}

variable "data_volumes" {
  type = list(object({
    ebs_volume_id     = string
    availability_zone = string
  }))
  description = "List of EBS volumes"
}

variable "mongodb_version" {
  type        = string
  description = "MongoDB version"
  default     = "4.2"
}

variable "tags" {
  type        = map(string)
  description = "Tags"
  default     = {}
}

variable "keypair_name" {
  type        = string
  description = "Keypair name"
  default     = "mongo-publicKey"
}

variable "public_key" {
  type        = string
  description = "Public key file path"
}

variable "private_key" {
  type        = string
  description = "Private key file path"
}

variable "bastion_host" {
  type        = string
  description = "Bastion host IP"
  default     = ""
}

variable "ami_owners" {
  type        = list(string)
  description = "AMI owners filter"
  default     = ["self", "amazon", "aws-marketplace"]
}

variable "ssh_user" {
  type        = string
  description = "SSH user name"
}

variable "bastion_user" {
  type        = string
  description = "bastion SSH user name"
  default     = ""
}

variable "replicaset_name" {
  type        = string
  description = "MongoDB ReplicaSet Name"
  default     = ""
}

variable "replica_count" {
  type        = number
  description = "Number of Replica nodes"
  default     = 1
}
variable "security_authorization" {
  type        = string
  description = "Authorization"
  default     = "disabled"
}
variable "admin_user_name" {
  type        = string
  description = "Admin user name"
  default     = "admin"
}
variable "admin_user_password" {
  type        = string
  description = "Admin user password"
  default     = "passw0rd"
}
variable "root_user_name" {
  type        = string
  description = "Root user name"
  default     = "root"
}
variable "root_user_password" {
  type        = string
  description = "Root user password"
  default     = "passw0rd"
}
variable "backup_user_name" {
  type        = string
  description = "Backup user name"
  default     = "root"
}
variable "backup_user_password" {
  type        = string
  description = "Backup user password"
  default     = "passw0rd"
}
variable "keyfile" {
  type        = string
  description = "keyfile path"
  default     = "./keyfile"
}
