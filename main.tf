locals {
  ami_id             = var.ami_id == "" ? data.aws_ami.ami.id : var.ami_id
  public_key_name    = var.keypair_name
  device_name        = "/dev/xvdh"
  ansible_host_group = ["db-mongodb"]
  replica_count      = var.replica_count < 1 ? 1 : var.replica_count
}

data "aws_vpc" "selected_vpc" {
  id = var.vpc_id
}

resource "aws_key_pair" "mongo_keypair" {
  key_name   = local.public_key_name
  public_key = var.public_key
}

resource "aws_instance" "mongo_server" {
  count                       = local.replica_count
  ami                         = local.ami_id
  instance_type               = var.instance_type
  subnet_id                   = var.subnet_id
  vpc_security_group_ids      = [aws_security_group.sg_mongodb.id]
  key_name                    = aws_key_pair.mongo_keypair.key_name
  availability_zone           = var.data_volumes[count.index].availability_zone
  associate_public_ip_address = true
  tags                        = var.tags

  connection {
    host         = var.bastion_host == "" ? self.public_ip : self.private_ip
    type         = "ssh"
    user         = var.ssh_user
    private_key  = var.private_key
    bastion_host = var.bastion_host
    agent        = true
  }

  provisioner "file" {
    source      = "${path.module}/provisioning/wait-for-cloud-init.sh"
    destination = "/tmp/wait-for-cloud-init.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo ln -s /usr/bin/python3 /usr/bin/python",
      "chmod +x /tmp/wait-for-cloud-init.sh",
      "/tmp/wait-for-cloud-init.sh",
    ]
  }
}

resource "aws_volume_attachment" "mongo-data-vol-attachment" {
  count       = local.replica_count
  device_name = local.device_name
  volume_id   = var.data_volumes[count.index].ebs_volume_id
  instance_id = aws_instance.mongo_server[count.index].id

  skip_destroy = true

  connection {
    host         = var.bastion_host == "" ? aws_instance.mongo_server[count.index].public_ip : aws_instance.mongo_server[count.index].private_ip
    type         = "ssh"
    user         = var.ssh_user
    private_key  = var.private_key
    bastion_host = var.bastion_host
    agent        = true
  }

  provisioner "file" {
    source      = "${path.module}/provisioning/mount-data-volume.sh"
    destination = "/tmp/mount-data-volume.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/mount-data-volume.sh",
      "/tmp/mount-data-volume.sh",
    ]
  }

  provisioner "ansible" {
    plays {
      playbook {
        file_path = "${path.module}/provisioning/playbook.yaml"
      }
      extra_vars = {
        mongodb_version             = var.mongodb_version
        mongodb_replication_replset = var.replicaset_name
      }
      groups = local.ansible_host_group
    }
  }
}

resource "null_resource" "replicaset_initialization" {
  depends_on = [aws_volume_attachment.mongo-data-vol-attachment]

  provisioner "file" {
    content = templatefile("${path.module}/provisioning/init-replicaset.js.tmpl", {
      replicaSetName = var.replicaset_name
      ip_addrs       = var.bastion_host == "" ? aws_instance.mongo_server.*.public_ip : aws_instance.mongo_server.*.private_ip

    })
    destination = "/tmp/init-replicaset.js"
  }

  provisioner "remote-exec" {
    inline = [
      "mongo 127.0.0.1:27017/admin /tmp/init-replicaset.js",
    ]
  }

  connection {
    host         = var.bastion_host == "" ? aws_instance.mongo_server[0].public_ip : aws_instance.mongo_server[0].private_ip
    type         = "ssh"
    user         = var.ssh_user
    private_key  = var.private_key
    bastion_host = var.bastion_host
    agent        = true
  }
}
