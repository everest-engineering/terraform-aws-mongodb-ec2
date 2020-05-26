output "mongo_server_private_ip" {
  description = "MongoDB Primary Private IP Address"
  value = aws_instance.mongo_server[0].private_ip
}

output "mongo_server_public_ip" {
  description = "MongoDB Primary Public IP Address"
  value = aws_instance.mongo_server[0].public_ip
}

output "mongo_replica_public_ip" {
  description = "MongoDB Replica IP Addresses"
  value = aws_instance.mongo_server.*.public_ip
}
