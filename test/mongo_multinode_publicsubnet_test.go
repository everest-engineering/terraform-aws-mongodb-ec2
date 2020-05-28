package test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestMongoDBMultiNodeInPublicSubnet(t *testing.T) {
	replicaCount := 2
	prerequisiteTerraformOptions := setupMongoInPublicSubnetPrerequisites(t, replicaCount)
	dataVolumes := getDataVolumes(t, prerequisiteTerraformOptions)
	vars := map[string]interface{}{
		"region":        "us-east-1",
		"keypair_name":  "mongo_keypair_tmp",
		"replica_count": replicaCount,
		"data_volumes":  dataVolumes,
	}
	terraformOptions := setupMongoInPublicSubnet(t, vars)
	mongoPublicIPAddresses := terraform.OutputList(t, terraformOptions, "mongo_replica_ip_address")

	primaryIP := mongoPublicIPAddresses[0]
	fmt.Println("MongoDB Primary IP: ", primaryIP)
	checkMongoConnectivity(t, primaryIP)

	secondaryIP := mongoPublicIPAddresses[1]
	fmt.Println("MongoDB Secondary IP: ", secondaryIP)
	checkMongoConnectivity(t, secondaryIP)

	checkMongoReplication(t, primaryIP, secondaryIP)

	defer terraform.Destroy(t, prerequisiteTerraformOptions)
	defer terraform.Destroy(t, terraformOptions)
}

func checkMongoReplication(t *testing.T, primaryIp string, secondaryIp string) {
	primaryURL := "mongodb://" + primaryIp + ":27017"

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(primaryURL))
	assert.Nil(t, err)

	database := client.Database("admin")

	var result bson.M
	cmd := bson.M{
		"replSetGetStatus": 1,
	}
	err = database.RunCommand(ctx, cmd).Decode(&result)
	assert.Nil(t, err)

	replMembers := []interface{}(result["members"].(primitive.A))
	assert.Equal(t, len(replMembers), 2)

	secondaryReplStatus := replMembers[1].(primitive.M)

	assert.Equal(t, secondaryReplStatus["name"], secondaryIp+":27017")
	assert.Equal(t, secondaryReplStatus["stateStr"], "SECONDARY")
	assert.Equal(t, secondaryReplStatus["health"], 1.0)
}
