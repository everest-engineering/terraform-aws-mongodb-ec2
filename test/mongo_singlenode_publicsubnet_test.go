package test

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var replicaCount = 1

func TestMongoDBSingleNodeInPublicSubnet(t *testing.T) {
	prerequisiteTerraformOptions := setupMongoInPublicSubnetPrerequisites(t)
	dataVolumes := getDataVolumes(t, prerequisiteTerraformOptions)
	vars := map[string]interface{}{
		"region":        "us-east-1",
		"keypair_name":  "mongo_keypair_tmp",
		"replica_count": replicaCount,
		"data_volumes":  dataVolumes,
	}
	terraformOptions := setupMongoInPublicSubnet(t, vars)

	defer terraform.Destroy(t, prerequisiteTerraformOptions)
	defer terraform.Destroy(t, terraformOptions)

	mongodbPublicIp := terraform.Output(t, terraformOptions, "mongo_server_ip_address")
	fmt.Println("mongodb public ip: ", mongodbPublicIp)
	checkMongoConnectivity(t, mongodbPublicIp)
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

func checkMongoConnectivity(t *testing.T, mongodbIp string) {
	mongodbConnectUrl := "mongodb://" + mongodbIp + ":27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbConnectUrl))
	assert.Nil(t, err)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	assert.Nil(t, err)

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	assert.Nil(t, err)
}

func getDataVolumes(t *testing.T, terraformOptions *terraform.Options) []map[string]interface{} {
	ebsVolumeIds := terraform.OutputList(t, terraformOptions, "ebs-vol-id")
	azs := terraform.OutputList(t, terraformOptions, "availability_zone")
	var dataVolumes []map[string]interface{}
	for i := 0; i< len(ebsVolumeIds); i++ {
		dataVolumes = append(dataVolumes, map[string]interface{}{
			"ebs_volume_id": ebsVolumeIds[i],
			"availability_zone": azs[i],
		})
	}
	fmt.Println("dataVolumes: ", dataVolumes)
	return dataVolumes
}

func setupMongoInPublicSubnetPrerequisites(t *testing.T) *terraform.Options {
    terraformOptions := &terraform.Options{
        TerraformDir: "../examples/mongodb-in-public-subnet/prerequisite",
        Vars: map[string]interface{}{
        	"volume_count": replicaCount,
        },
    }
    terraform.InitAndApply(t, terraformOptions)
    return terraformOptions
}
