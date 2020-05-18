output "mongo_server_private_ip" {
  value = aws_instance.mongo_server[0].private_ip
}

output "mongo_server_public_ip" {
  value = aws_instance.mongo_server[0].public_ip
}
