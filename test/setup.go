package test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"testing"
)

func setupMongoInPublicSubnetPrerequisites(t *testing.T, replicaCount int) *terraform.Options {
	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/mongodb-in-public-subnet/prerequisite",
		Vars: map[string]interface{}{
			"volume_count": replicaCount,
		},
	}
	terraform.InitAndApply(t, terraformOptions)
	return terraformOptions
}

func setupMongoInPublicSubnet(t *testing.T, vars map[string]interface{}) *terraform.Options {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/mongodb-in-public-subnet",
		// Variables to pass to our Terraform code using -var options
		Vars: vars,
	}

	terraform.InitAndApply(t, terraformOptions)
	return terraformOptions
}

func getDataVolumes(t *testing.T, terraformOptions *terraform.Options) []map[string]interface{} {
	ebsVolumeIds := terraform.OutputList(t, terraformOptions, "ebs-vol-id")
	azs := terraform.OutputList(t, terraformOptions, "availability_zone")
	var dataVolumes []map[string]interface{}
	for i := 0; i < len(ebsVolumeIds); i++ {
		dataVolumes = append(dataVolumes, map[string]interface{}{
			"ebs_volume_id":     ebsVolumeIds[i],
			"availability_zone": azs[i],
		})
	}
	fmt.Println("dataVolumes: ", dataVolumes)
	return dataVolumes
}
