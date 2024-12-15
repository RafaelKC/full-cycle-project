package main

import (
	"context"
	"fmt"
	"github.com/RafaelKC/full-cycle-project/golang-simulador/internal"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	mongoStr := os.Getenv("MONGODB_URL")
	mongoConnection, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoStr))
	if err != nil {
		panic(err)
	}

	freightService := internal.NewFreightService()
	routeService := internal.NewRouteService(mongoConnection, freightService)

	kafkaAddr := os.Getenv("KAFKA_URL")
	freightWriter := &kafka.Writer{
		Addr:     kafka.TCP(kafkaAddr),
		Topic:    "freight",
		Balancer: &kafka.LeastBytes{},
	}
	simulatorWriter := &kafka.Writer{
		Addr:     kafka.TCP(kafkaAddr),
		Topic:    "simulator",
		Balancer: &kafka.LeastBytes{},
	}

	hub := internal.NewEventHub(mongoConnection, routeService, freightWriter, simulatorWriter)

	routeReader := *kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaAddr},
		Topic:   "route",
		GroupID: "simulator",
	})

	fmt.Println("Starting simulator")
	for {
		m, err := routeReader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		go func(msg []byte) {
			err = hub.HandleEvent(m.Value)
			if err != nil {
				return
			}
		}(m.Value)
	}
}
