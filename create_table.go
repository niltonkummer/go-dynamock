package dynamock

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Name - method for set Name expectation
func (e *CreateTableExpectation) Name(table string) *CreateTableExpectation {
	e.table = &table
	return e
}

// KeySchema - method for set KeySchema expectation
func (e *CreateTableExpectation) KeySchema(keySchema []*types.KeySchemaElement) *CreateTableExpectation {
	e.keySchema = keySchema
	return e
}

// WillReturns - method for set desired result
func (e *CreateTableExpectation) WillReturns(res dynamodb.CreateTableOutput) *CreateTableExpectation {
	e.output = &res
	return e
}

// CreateTable - this func will be invoked when test running matching expectation with actual input
func (e *MockDynamoDB) CreateTable(ctx context.Context, input *dynamodb.CreateTableInput, opts ...func(*dynamodb.Options))  (*dynamodb.CreateTableOutput, error) {
	if len(e.dynaMock.CreateTableExpect) > 0 {
		x := e.dynaMock.CreateTableExpect[0] //get first element of expectation

		if x.table != nil {
			if *x.table != *input.TableName {
				return &dynamodb.CreateTableOutput{}, fmt.Errorf("Expect table %s but found table %s", *x.table, *input.TableName)
			}
		}

		if x.keySchema != nil {
			if !reflect.DeepEqual(x.keySchema, input.KeySchema) {
				return &dynamodb.CreateTableOutput{}, fmt.Errorf("Expect keySchema %+v but found keySchema %+v", x.keySchema, input.KeySchema)
			}
		}

		// delete first element of expectation
		e.dynaMock.CreateTableExpect = append(e.dynaMock.CreateTableExpect[:0], e.dynaMock.CreateTableExpect[1:]...)

		return x.output, nil
	}

	return &dynamodb.CreateTableOutput{}, fmt.Errorf("Create Table Expectation Not Found")
}
