package infrastructure

import (
	"fmt"
	"strconv"

	"github.com/ZOLUXERO/gotest/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBUserRepository struct {
	Svc       *dynamodb.DynamoDB
	tableName string
}

func (r *DynamoDBUserRepository) GetByID(userID string) (*core.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("users"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userID),
			},
		},
	}

	result, err := r.Svc.GetItem(input)
	if err != nil {
		return nil, err
	}

	item := result.Item
	if item == nil {
		return nil, nil
	}
	points, err := strconv.Atoi(*item["Points"].N)
	if err != nil {
		return nil, nil
	}
	user := &core.User{
		ID:     *item["ID"].S,
		Points: points,
	}
	return user, nil
}

func (r *DynamoDBUserRepository) Save(user *core.User) error {
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

	_, err := r.Svc.PutItem(input)
	return err
}

func (r *DynamoDBUserRepository) AddPoints(user *core.User, points int64) error {

	result, err := r.Svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("users"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(user.ID)},
		},
	})
	if err != nil {
		return err
	}
	item := result.Item
	currentPoints, err := strconv.Atoi(*item["Points"].N)
	if err != nil {
		return err
	}

	//  Sumar puntos
	newPoints := currentPoints + int(points)

	// el error (ValidationException: The number of conditions on the keys is invalid) proviene de aca
	_, err = r.Svc.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String("users"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {S: aws.String(user.ID)},
		},
		UpdateExpression: aws.String("set Points = :p"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p": {N: aws.String(strconv.Itoa(newPoints))},
		},
	})
	return err
}

func (r *DynamoDBUserRepository) SubtractPoints(user *core.User, points int64) error {
	pointsToSubtract := int64(points)
	currentPoints := user.Points

	if currentPoints < int(pointsToSubtract) {
		return fmt.Errorf("no tiene suficientes puntos")
	}
	newPoints := currentPoints - int(pointsToSubtract)
	updateExpression := "SET Points = :p"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":p": {
			N: aws.String(strconv.FormatInt(int64(newPoints), 10)),
		},
	}

	// construya lo que va a actualizar
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(user.ID),
			},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              aws.String("ACTUALIZADO"),
	}

	// actualize
	_, err := r.Svc.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("error updating user points: %w", err)
	}
	return nil
}

// Problema al instanciar region local??????????
/* func NewDynamoDBUserRepository() (*DynamoDBUserRepository, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("local"),
		Endpoint:    aws.String("http://localhost:8000"),
		Credentials: credentials.NewStaticCredentials("ACCESS_KEY", "SECRET_KEY", ""),
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	return &DynamoDBUserRepository{
		Svc: dynamodb.New(sess),
	}, nil
} */

func NewDynamoDBUserRepository(svc *dynamodb.DynamoDB, tableName string) *DynamoDBUserRepository {
	return &DynamoDBUserRepository{
		Svc:       svc,
		tableName: tableName,
	}
}
