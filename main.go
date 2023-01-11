package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gotest/infrastructure"
)

func main() {
	// Create an instance of the router
	router := gin.Default()

	// Create instances of the services and repositories
	userRepo := &infrastructure.DynamoDBUserRepository{
		svc: dynamodb.New(session.New()),
	}
	kafkaEmitter := &infrastructure.KafkaEventEmitter{
		producer: &kafka.Producer{},
	}
	userService := &infrastructure.UserService{
		repo:         userRepo,
		kafkaEmitter: kafkaEmitter,
	}

	// Define the API endpoints and handlers
	router.PUT("/user/points/accumulate/:userId/:points/:type", func(c *gin.Context) {
		userId := c.Param("userId")
		points, _ := strconv.Atoi(c.Param("points"))
		pointType := c.Param("type")

		err := userService.AccumulatePoints(userId, points, pointType)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Points accumulated"})
	})
	router.PUT("/user/points/redeem/:userId/:points", func(c *gin.Context) {
		userId := c.Param("userId")
		points, _ := strconv.Atoi(c.Param("points"))

		err := userService.RedeemPoints(userId, points)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Points redeemed"})
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
	// Start the server
	router.Run(":8080")
}
