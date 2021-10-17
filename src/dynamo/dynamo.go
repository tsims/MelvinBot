package dynamo

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// NewDynamoSession is designed to run inside an EC2 environment with IAM role
// already configured with access to Dynamo

type DynamoClient struct {
	client *dynamodb.DynamoDB
	lock   *sync.Mutex
}

type Pin struct {
	Number      int
	Author      string
	Message     string
	Time        time.Time
	MessageLink string
}

func NewDynamoSession() DynamoClient {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("failed to create dynamo session", err)
	}

	return DynamoClient{
		client: dynamodb.New(sess),
		lock:   &sync.Mutex{},
	}
}

const (
	tableName string = "MelvinPins"
	partKey   string = "PinNumber"
)

func (d *DynamoClient) GetPin(number int) (Pin, error) {

	pinNum := strconv.Itoa(number)
	input := dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			partKey: {
				N: aws.String(pinNum),
			},
		},
	}

	pinItem, err := d.client.GetItem(&input)
	if err != nil {
		return Pin{}, err
	}

	if pinItem.Item == nil {
		return Pin{}, fmt.Errorf("Pin does not exist at number %d", number)
	}

	pin := Pin{}
	err = dynamodbattribute.UnmarshalMap(pinItem.Item, &pin)
	if err != nil {
		return Pin{}, err
	}

	return pin, nil
}
