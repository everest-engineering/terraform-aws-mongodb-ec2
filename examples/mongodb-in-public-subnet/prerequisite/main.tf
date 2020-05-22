provider "aws" {
  region = "us-east-1"
}

resource "aws_ebs_volume" "mongo-data-vol" {
  count             = var.volume_count
  availability_zone = var.availability_zone
  type              = "gp2"
  size              = var.volume_size

  tags = {
    Name        = "mongo-data-ebs-volume"
    Environment = var.environment_tag
  }
}

variable "volume_count" {
  type        = number
  description = "Number of volumes to create"
  default     = 1
}

variable "availability_zone" {
  type        = string
  description = "Availability zone"
  default     = "us-east-1a"
}

variable "volume_size" {
  type        = string
  description = "Size of the DB storage volume."
  default     = "10"
}

variable "environment_tag" {
  type        = string
  description = "Environment tag"
  default     = "Production"
}

output "ebs-vol-id" {
  value = aws_ebs_volume.mongo-data-vol.*.id
}

output "availability_zone" {
  value = aws_ebs_volume.mongo-data-vol.*.availability_zone
}