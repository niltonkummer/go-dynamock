package examples

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/niltonkummer/dynamo/dynamodbiface"
)

// MyDynamo struct hold dynamodb connection
type MyDynamo struct {
	Db dynamodbiface.DynamoDBAPI
}

// Dyna - object from MyDynamo
var Dyna *MyDynamo

// ConfigureDynamoDB - init func for open connection to aws dynamodb
func ConfigureDynamoDB() {
	Dyna = new(MyDynamo)
	awsConfig, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-2"))
	svc := dynamodb.NewFromConfig(awsConfig)
	Dyna.Db = dynamodbiface.DynamoDBAPI(svc)
}

// GetName - example func using GetItem method
func GetName(id string) (string, error) {
	parameter := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: id},
		},
		TableName: aws.String("employee"),
	}

	response, err := Dyna.Db.GetItem(context.TODO(), parameter)
	if err != nil {
		return "", err
	}

	name := response.Item["name"].(*types.AttributeValueMemberS).Value
	return name, nil
}

// GetName - example func using GetItem method
func GetTransactGetItems(id string) error {
	parameter := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName: aws.String("my_table"),
				},
			},
		},
	}

	_, err := Dyna.Db.TransactWriteItems(context.TODO(), parameter)

	if err != nil {
		fmt.Print(err.Error())
		return err
	}

	return nil
}
