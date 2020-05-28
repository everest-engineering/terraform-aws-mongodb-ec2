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

func TestMongoDBSingleNodeInPublicSubnet(t *testing.T) {
	var replicaCount = 1
	prerequisiteTerraformOptions := setupMongoInPublicSubnetPrerequisites(t, replicaCount)
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
