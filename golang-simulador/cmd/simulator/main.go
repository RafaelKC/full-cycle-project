package main

import (
	"context"
	"fmt"
	"github.com/RafaelKC/full-cycle-project/golang-simulador/internal"
	"github.com/joho/godotenv"
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

	routeCreatedEvent := internal.NewRouteCreatedEvent(
		"1",
		100,
		[]internal.Direction{
			{Lat: 0, Lng: 0},
			{Lat: 10, Lng: 10},
		},
	)

	fmt.Println(internal.RouteCreatedHandler(routeCreatedEvent, routeService))
}
