package email

import (
	"bytes"
	"html/template"
	"time"
)

type Sender struct {
	client *SMTPClient
	tmpl   *template.Template
}

func NewSender(client *SMTPClient, templateDir string) (*Sender, error) {
	tmpl, err := template.ParseGlob(templateDir + "/*.html")
	if err != nil {
		return nil, err
	}
	return &Sender{
		client: client,
		tmpl:   tmpl,
	}, nil
}

func sendWithRetry(send func() error) error {
	delays := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}
	var lastErr error
	for _, d := range delays {
		if lastErr = send(); lastErr == nil {
			return nil
		}
		time.Sleep(d)
	}
	return lastErr
}

func (s *Sender) SendBookingConfirmation(email, movie, seat, time string) error {
	var buf bytes.Buffer
	data := map[string]interface{}{
		"Movie": movie,
		"Seat":  seat,
		"Time":  time,
	}
	if err := s.tmpl.ExecuteTemplate(&buf, "booking_confirmation.html", data); err != nil {
		return err
	}

	return sendWithRetry(func() error {
		return s.client.Send(email, "Booking Confirmation", buf.String())
	})
}

func (s *Sender) SendCancellation(email, bookingID, date string) error {
	var buf bytes.Buffer
	data := map[string]interface{}{
		"BookingID": bookingID,
		"Date":      date,
	}
	if err := s.tmpl.ExecuteTemplate(&buf, "booking_cancellation.html", data); err != nil {
		return err
	}

	return sendWithRetry(func() error {
		return s.client.Send(email, "Booking Cancellation", buf.String())
	})
}

func (s *Sender) SendPaymentReceipt(email string, amount float64, paymentID string) error {
	var buf bytes.Buffer
	data := map[string]interface{}{
		"Amount":    amount,
		"PaymentID": paymentID,
	}
	if err := s.tmpl.ExecuteTemplate(&buf, "payment_receipt.html", data); err != nil {
		return err
	}

	return sendWithRetry(func() error {
		return s.client.Send(email, "Payment Receipt", buf.String())
	})
}
