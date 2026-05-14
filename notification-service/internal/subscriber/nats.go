package subscriber

import (
	"encoding/json"
	"log"

	"notification-service/internal/email"
	"github.com/nats-io/nats.go"
)

type NATSSubscriber struct {
	conn   *nats.Conn
	sender *email.Sender
}

func NewNATSSubscriber(conn *nats.Conn, sender *email.Sender) *NATSSubscriber {
	return &NATSSubscriber{
		conn:   conn,
		sender: sender,
	}
}

type BookingCreatedEvent struct {
	BookingID string `json:"booking_id"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Movie     string `json:"movie"`
	Seat      string `json:"seat"`
	Time      string `json:"time"`
}

type BookingCancelledEvent struct {
	BookingID string `json:"booking_id"`
	Email     string `json:"email"` // Expecting email and date for cancellation
	Date      string `json:"date"`
}

type PaymentConfirmedEvent struct {
	PaymentID string  `json:"payment_id"`
	BookingID string  `json:"booking_id"`
	Amount    float64 `json:"amount"`
	Email     string  `json:"email"`
}

func (s *NATSSubscriber) Start() error {
	_, err := s.conn.Subscribe("booking.created", func(msg *nats.Msg) {
		var event BookingCreatedEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Error decoding booking.created event: %v", err)
			return
		}
		if err := s.sender.SendBookingConfirmation(event.Email, event.Movie, event.Seat, event.Time); err != nil {
			log.Printf("Failed to send booking confirmation email: %v", err)
		} else {
			log.Printf("Sent booking confirmation to %s", event.Email)
		}
	})
	if err != nil {
		return err
	}

	_, err = s.conn.Subscribe("booking.cancelled", func(msg *nats.Msg) {
		var event BookingCancelledEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Error decoding booking.cancelled event: %v", err)
			return
		}
		if err := s.sender.SendCancellation(event.Email, event.BookingID, event.Date); err != nil {
			log.Printf("Failed to send booking cancellation email: %v", err)
		} else {
			log.Printf("Sent booking cancellation to %s", event.Email)
		}
	})
	if err != nil {
		return err
	}

	_, err = s.conn.Subscribe("payment.confirmed", func(msg *nats.Msg) {
		var event PaymentConfirmedEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Error decoding payment.confirmed event: %v", err)
			return
		}
		if err := s.sender.SendPaymentReceipt(event.Email, event.Amount, event.PaymentID); err != nil {
			log.Printf("Failed to send payment receipt email: %v", err)
		} else {
			log.Printf("Sent payment receipt to %s", event.Email)
		}
	})
	return err
}
