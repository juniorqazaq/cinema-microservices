package publisher

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type NATSPublisher struct {
	conn *nats.Conn
}

func NewNATSPublisher(conn *nats.Conn) *NATSPublisher {
	return &NATSPublisher{conn: conn}
}

type BookingCreatedEvent struct {
	BookingID string `json:"booking_id"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Movie     string `json:"movie"`
	Seat      string `json:"seat"`
	Time      string `json:"time"`
}

func (p *NATSPublisher) PublishBookingCreated(ctx context.Context, event *BookingCreatedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.conn.Publish("booking.created", data)
}

type BookingCancelledEvent struct {
	BookingID string `json:"booking_id"`
}

func (p *NATSPublisher) PublishBookingCancelled(ctx context.Context, event *BookingCancelledEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.conn.Publish("booking.cancelled", data)
}

type PaymentConfirmedEvent struct {
	PaymentID string  `json:"payment_id"`
	BookingID string  `json:"booking_id"`
	Amount    float64 `json:"amount"`
	Email     string  `json:"email"`
}

func (p *NATSPublisher) PublishPaymentConfirmed(ctx context.Context, event *PaymentConfirmedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.conn.Publish("payment.confirmed", data)
}
