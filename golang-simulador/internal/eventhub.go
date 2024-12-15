package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type EventHub struct {
	routeService    *RouteService
	mongoClient     *mongo.Client
	chDriverMoved   chan *DriverMovedEvent
	freightWriter   *kafka.Writer
	simulatorWriter *kafka.Writer
}

func (eh *EventHub) HandleEvent(msg []byte) error {
	var baseEvent struct {
		EventName string `json:"event_name"`
	}
	err := json.Unmarshal(msg, &baseEvent)
	if err != nil {
		return fmt.Errorf("unmarshal event: %w", err)
	}

	switch baseEvent.EventName {
	case "RouteCreated":
		var event RouteCreatedEvent
		err = json.Unmarshal(msg, &event)
		if err != nil {
			return fmt.Errorf("unmarshal event: %w", err)
		}
		return eh.handleRouteCreated(event)

	case "DeliveryStarted":
		var event DeliveryStartedEvent
		err = json.Unmarshal(msg, &event)
		if err != nil {
			return fmt.Errorf("unmarshal event: %w", err)
		}
		return eh.handleDeliveryStarted(event)

	default:
		return errors.New("unknown event")
	}
}

func (eh *EventHub) handleRouteCreated(event RouteCreatedEvent) error {
	freightCalculatedEvent, err := RouteCreatedHandler(&event, eh.routeService)
	if err != nil {
		return err
	}
	value, err := json.Marshal(freightCalculatedEvent)
	if err != nil {
		return err
	}

	err = eh.freightWriter.WriteMessages(context.Background(), kafka.Message{
		Value: value,
		Key:   []byte(freightCalculatedEvent.RouteId),
	})
	if err != nil {
		return err
	}
	return nil
}

func (eh *EventHub) handleDeliveryStarted(event DeliveryStartedEvent) error {
	err := DeliveryStatedHandler(&event, eh.routeService, eh.chDriverMoved)
	if err != nil {
		return err
	}
	go eh.sendDirections()
	return nil
}

func (eh *EventHub) sendDirections() {
	for {
		select {
		case movedEvent := <-eh.chDriverMoved:
			value, err := json.Marshal(movedEvent)
			if err != nil {
				return
			}
			err = eh.simulatorWriter.WriteMessages(context.Background(), kafka.Message{
				Value: value,
				Key:   []byte(movedEvent.RouteId),
			})
			if err != nil {
				return
			}
		case <-time.After(500 * time.Millisecond):
			return
		}
	}
}

func NewEventHub(mongoClient *mongo.Client, routeService *RouteService, freightWriter *kafka.Writer, simulatorWriter *kafka.Writer) *EventHub {
	return &EventHub{
		routeService:    routeService,
		mongoClient:     mongoClient,
		freightWriter:   freightWriter,
		simulatorWriter: simulatorWriter,
		chDriverMoved:   make(chan *DriverMovedEvent),
	}
}
