package internal

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
)

type Direction struct {
	Lat float64 `bson:"lat" json:"lat"`
	Lng float64 `bson:"lng" json:"lng"`
}

type Route struct {
	Id           string      `bson:"_id" json:"id"`
	Distance     int         `bson:"distance" json:"distance"`
	FreightPrice float64     `bson:"freight_price" json:"freight_rice"`
	Directions   []Direction `bson:"directions" json:"directions"`
}

func NewRoute(id string, distance int, directions []Direction) *Route {
	return &Route{
		Id:         id,
		Distance:   distance,
		Directions: directions,
	}
}

type FreightService struct{}

func (fs *FreightService) Calculate(distance int) float64 {
	return math.Floor((float64(distance)*0.15+0.3)*100) / 100
}

func NewFreightService() *FreightService {
	return &FreightService{}
}

type RouteService struct {
	mongo          *mongo.Client
	freightService *FreightService
}

func (rs *RouteService) CreateRoute(route *Route) (*Route, error) {
	route.FreightPrice = rs.freightService.Calculate(route.Distance)

	update := bson.M{
		"$set": bson.M{
			"distance":      route.Distance,
			"directions":    route.Directions,
			"freight_price": route.FreightPrice,
		},
	}

	filter := bson.M{"_id": route.Id}
	opts := options.Update().SetUpsert(true)

	_, err := rs.mongo.Database("go").Collection("routes").UpdateOne(nil, filter, update, opts)
	if err != nil {
		return nil, err
	}

	return route, nil
}

func (rs *RouteService) GetRoute(routeId string) (Route, error) {
	var route Route
	filter := bson.M{"_id": routeId}
	err := rs.mongo.Database("go").Collection("routes").FindOne(nil, filter).Decode(&route)
	if err != nil {
		return Route{}, err
	}

	return route, nil
}

func NewRouteService(mongo *mongo.Client, freightService *FreightService) *RouteService {
	return &RouteService{
		mongo:          mongo,
		freightService: freightService,
	}
}
