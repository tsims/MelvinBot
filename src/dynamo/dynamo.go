package dynamo

import (
	"fmt"
	"sync"

	stats "MelvinBot/src/stats"

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

func NewDynamoSession() (DynamoClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	})
	if err != nil {
		return DynamoClient{}, err
	}

	ddb := dynamodb.New(sess)

	return DynamoClient{
		client: ddb,
		lock:   &sync.Mutex{},
	}, nil
}

const (
	statsTable string = "stats"
	partKey    string = "guild"
	statsKey   string = "stats"
)

func (d *DynamoClient) GetStatsOnAllGuilds() error {

	items, err := d.client.Scan(&dynamodb.ScanInput{
		TableName: aws.String(statsTable)})

	if err != nil {
		return fmt.Errorf("couldnt scan dynamo: %v", err)
	}

	statsHolder := map[string]*stats.Stats{}

	for _, item := range items.Items {
		var guild string
		err := dynamodbattribute.Unmarshal(item[partKey], &guild)
		if err != nil {
			return fmt.Errorf("couldnt unmarshal struct %s: %v", partKey, err)
		}
		var statsUnmarshal map[string]int
		err = dynamodbattribute.Unmarshal(item[statsKey], &statsUnmarshal)
		if err != nil {
			return fmt.Errorf("couldnt unmarshal struct %s: %v", statsKey, err)
		}
		statsStruct := stats.Stats{
			StatMap: statsUnmarshal,
			Lock:    &sync.Mutex{},
		}
		statsHolder[guild] = &statsStruct
	}

	if len(statsHolder) != 0 {
		stats.StatsPerGuild = statsHolder
	}

	return nil
}

func (d *DynamoClient) PutStatsOnAllGuilds() error {

	input := map[string]*dynamodb.AttributeValue{}
	for guild, stats := range stats.StatsPerGuild {

		av, err := dynamodbattribute.Marshal(stats.StatMap)
		if err != nil {
			return fmt.Errorf("couldnt marshal struct: %v", err)
		}
		guildStr, err := dynamodbattribute.Marshal(guild)
		if err != nil {
			return fmt.Errorf("couldnt marshal struct: %v", err)
		}
		input[partKey] = guildStr
		input[statsKey] = av

		_, err = d.client.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(statsTable),
			Item:      input,
		})
		if err != nil {
			return fmt.Errorf("couldnt put dynamo: %v", err)
		}
	}

	return nil
}
