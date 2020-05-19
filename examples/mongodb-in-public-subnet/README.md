# Provision MongoDB in Public Subnet
This example provision mongodb server in default VPC and in a public subnet.

1. Configure AWS Credentials as environment variables:

```shell script
export AWS_ACCESS_KEY_ID="ACCESS_KEY_HERE"
export AWS_SECRET_ACCESS_KEY="SECRET_ACCESS_KEY_HERE"
export AWS_DEFAULT_REGION="REGION_HERE"
```

2. If you don't have any existing EBS volume then create EBS volume using terraform script in prerequisite

```shell script
cd terraform-mongodb-provisioning-ec2/examples/mongodb-in-public-subnet/prerequisite
terraform init
terraform plan
terraform apply
```

3. Provision MongoDB:
Configure the existing/created `ebs_volume_id`(s) and a`vailability_zone`(s) as `data_volumes` variable in `mongodb-in-public-subnet/main.tf`.

```hcl-terraform
variable "data_volumes" {
  type = list(object({
    ebs_volume_id     = string
    availability_zone = string
  }))
  description = "List of EBS volumes"
  default     = [
    {
      ebs_volume_id = "EBS_VOLUME_ID_HERE"
      availability_zone = "AZ_HERE"
    },
    {
      ebs_volume_id = "EBS_VOLUME_ID_HERE"
      availability_zone = "AZ_HERE"
    }
  ]
}
```

Now you can provision MongoDB as follows:

```shell script
cd terraform-mongodb-provisioning-ec2/examples/mongodb-in-public-subnet
terraform init
terraform plan
terraform apply
```

4. Destroy:

```shell script
cd terraform-mongodb-provisioning-ec2/examples/mongodb-in-public-subnet
terraform destroy
```