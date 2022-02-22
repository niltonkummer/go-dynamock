package dynamock

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Table - method for set Table expectation
func (e *TransactWriteItemsExpectation) Table(table string) *TransactWriteItemsExpectation {
	e.table = &table
	return e
}

// WithItems - method for set Items expectation
func (e *TransactWriteItemsExpectation) WithItems(items []*types.TransactWriteItem) *TransactWriteItemsExpectation {
	e.items = items
	return e
}

// WillReturns - method for set desired result
func (e *TransactWriteItemsExpectation) WillReturns(res dynamodb.TransactWriteItemsOutput) *TransactWriteItemsExpectation {
	e.output = &res
	return e
}

// TransactWriteItems - this func will be invoked when test running matching expectation with actual input
func (e *MockDynamoDB) TransactWriteItems(ctx context.Context, input *dynamodb.TransactWriteItemsInput, opts ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
	if len(e.dynaMock.TransactWriteItemsExpect) > 0 {
		x := e.dynaMock.TransactWriteItemsExpect[0] // get first element of expectation

		// compare lengths
		if len(x.items) != len(input.TransactItems) {
			return &dynamodb.TransactWriteItemsOutput{}, fmt.Errorf("Expect items %+v but found items %+v", x.items, input.TransactItems)
		}

		for i, item := range input.TransactItems {
			// comapre table name for each write transact item with `x.table`
			if x.table != nil {
				if (item.Update != nil && x.table != item.Update.TableName) ||
					(item.Put != nil && x.table != item.Put.TableName) ||
					(item.Delete != nil && x.table != item.Delete.TableName) {
					return &dynamodb.TransactWriteItemsOutput{}, fmt.Errorf("Expect table %s not found", *x.table)
				}
			}

			// compare transact write item - each item also contains the table name
			if x.items[i] != nil && !reflect.DeepEqual(x.items[i], item) {
				return &dynamodb.TransactWriteItemsOutput{}, fmt.Errorf("Expect item %+v at index %d but found item %+v", x.items[i], i, item)
			}
		}

		// delete first element of expectation
		e.dynaMock.TransactWriteItemsExpect = append(e.dynaMock.TransactWriteItemsExpect[:0],
			e.dynaMock.TransactWriteItemsExpect[1:]...)

		return x.output, nil
	}

	return &dynamodb.TransactWriteItemsOutput{}, fmt.Errorf("Transact Write Items Table Expectation Not Found")
}

func (e *MockDynamoDB) TransactWriteItemsWithContext(ctx context.Context, input *dynamodb.TransactWriteItemsInput, opts ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
	return e.TransactWriteItems(ctx, input, opts...)
}
