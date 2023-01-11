package infrastructure

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBUserRepository struct {
	svc *dynamodb.DynamoDB
}

func (r *DynamoDBUserRepository) GetByID(userID string) (*core.User, error) {
	// Prepare the input for the GetItem method
	input := &dynamodb.GetItemInput{
		TableName: aws.String("users"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userID),
			},
		},
	}

	// Call the GetItem method
	result, err := r.svc.GetItem(input)
	if err != nil {
		return nil, err
	}

	// Extract the user from the result
	item := result.Item
	if item == nil {
		return nil, nil
	}
	user := &core.User{
		ID:     *item["ID"].S,
		Points: *item["Points"].N,
	}
	return user, nil
}

func (r *DynamoDBUserRepository) Save(user *core.User) error {
	// Prepare the input for the PutItem method
	input := &dynamodb.PutItemInput{
		TableName: aws.String("users"),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(user.ID),
			},
			"Points": {
				N: aws.String(strconv.Itoa(user.Points)),
			},
		},
	}

	// Call the PutItem method
	_, err := r.svc.PutItem(input)
	return err
}

// NewDynamoDBUserRepository returns a new instance of DynamoDBUserRepository
func NewDynamoDBUserRepository() (*DynamoDBUserRepository, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost:8000"),
		Region:   aws.String("us-west-2"),
	})
	if err != nil {
		return nil, err
	}

	return &DynamoDBUserRepository{
		svc: dynamodb.New(sess),
	}, nil
}
