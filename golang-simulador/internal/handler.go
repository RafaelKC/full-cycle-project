package internal

import "time"

type RouteCreatedEvent struct {
	EventName  string      `json:"event_name"`
	RouteId    string      `json:"id"`
	Distance   int         `json:"distance"`
	Directions []Direction `json:"directions"`
}

func NewRouteCreatedEvent(routeId string, distance int, directions []Direction) *RouteCreatedEvent {
	return &RouteCreatedEvent{
		EventName:  "RouteCreated",
		RouteId:    routeId,
		Distance:   distance,
		Directions: directions,
	}
}

type FreightCalculatedEvent struct {
	EventName string  `json:"event_name"`
	RouteId   string  `json:"route_id"`
	Amount    float64 `json:"amount"`
}

func NewFreightCalculatedEvent(routeId string, amount float64) *FreightCalculatedEvent {
	return &FreightCalculatedEvent{
		EventName: "FreightCalculated",
		RouteId:   routeId,
		Amount:    amount,
	}
}

type DeliveryStartedEvent struct {
	EventName string `json:"event_name"`
	RouteId   string `json:"route_id"`
}

func NewDeliveryStartedEvent(routeId string) *DeliveryStartedEvent {
	return &DeliveryStartedEvent{
		EventName: "DeliveryStated",
		RouteId:   routeId,
	}
}

type DriverMovedEvent struct {
	EventName string  `json:"event_name"`
	RouteId   string  `json:"route_id"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
}

func NewDiverMovedEvent(routeId string, lat float64, lng float64) *DriverMovedEvent {
	return &DriverMovedEvent{
		EventName: "DiverMoved",
		RouteId:   routeId,
		Lat:       lat,
		Lng:       lng,
	}
}

func RouteCreatedHandler(event *RouteCreatedEvent, rs *RouteService) (*FreightCalculatedEvent, error) {
	route := NewRoute(event.RouteId, event.Distance, event.Directions)
	routeCreate, err := rs.CreateRoute(route)
	if err != nil {
		return nil, err
	}

	freightCalculatedEvent := NewFreightCalculatedEvent(routeCreate.Id, routeCreate.FreightPrice)
	return freightCalculatedEvent, nil
}

func DeliveryStatedHandler(event *DeliveryStartedEvent, rs *RouteService, ch chan *DriverMovedEvent) error {
	route, err := rs.GetRoute(event.RouteId)
	if err != nil {
		return err
	}

	driverMovedEvent := NewDiverMovedEvent(event.RouteId, 0, 0)
	for _, direction := range route.Directions {
		driverMovedEvent.Lng = direction.Lng
		driverMovedEvent.Lat = direction.Lat
		time.Sleep(time.Second)
		ch <- driverMovedEvent
	}
	return nil
}
