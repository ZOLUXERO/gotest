package main

import (
	"log"
	"strconv"

	_ "github.com/ZOLUXERO/gotest/dynamoinitdb"
	"github.com/ZOLUXERO/gotest/infrastructure"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	/* db, err := infrastructure.NewDynamoDBUserRepository()
	if err != nil {
		log.Fatalf("failed to create session, %v", err)
	}

	// Instancias repositorios
	userRepo := &infrastructure.DynamoDBUserRepository{
		Svc: db,
	} */

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Endpoint:    aws.String("http://localhost:8000"),
		Credentials: credentials.NewStaticCredentials("ACCESS_KEY", "SECRET_KEY", ""),
	})

	if err != nil {
		log.Fatal("Error creating session: ", err)
	}

	db := dynamodb.New(sess)
	userRepo := infrastructure.NewDynamoDBUserRepository(db, "users")

	kafkaEmitter := &infrastructure.KafkaEventEmitter{
		Producer: &kafka.Producer{},
	}
	userService := &infrastructure.UserService{
		Repo:    userRepo,
		Emitter: kafkaEmitter,
	}

	// API endpoints
	router.PUT("/user/points/accumulate/:userId/:points/:type", func(c *gin.Context) {
		userId := c.Param("userId")
		points, _ := strconv.Atoi(c.Param("points"))
		pointType := c.Param("type")

		err := userService.AddPoints(userId, points, pointType)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"mensaje": "Puntos a;adidos"})
	})

	router.PUT("/user/points/redeem/:userId/:points", func(c *gin.Context) {
		userId := c.Param("userId")
		points, _ := strconv.Atoi(c.Param("points"))

		err := userService.SubtractPoints(userId, points, "new")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Puntos redimidos"})
	})

	router.GET("/user/points/:userId", func(c *gin.Context) {
		userId := c.Param("userId")
		points, err := userService.GetPoints(userId)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"points": points})
	})

	// Puerto en el que va iniciar el server
	router.Run(":8080")
}
