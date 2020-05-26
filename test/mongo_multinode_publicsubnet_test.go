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

	defer terraform.Destroy(t, prerequisiteTerraformOptions)
	defer terraform.Destroy(t, terraformOptions)

	primaryIP := mongoPublicIPAddresses[0]
	fmt.Println("MongoDB Primary IP: ", primaryIP)
	checkMongoConnectivity(t, primaryIP)

	secondaryIP := mongoPublicIPAddresses[1]
	fmt.Println("MongoDB Secondary IP: ", secondaryIP)
	checkMongoConnectivity(t, secondaryIP)
}
