package handler

import (
	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/services"
)

type PaymentHandlers struct {
	PaymentSchedule ResourceHandler
	Payment         ResourceHandler
}

func NewPaymentHandlers(s services.PaymentServices) PaymentHandlers {
	return PaymentHandlers{
		PaymentSchedule: NewCRUDHandler[models.PaymentSchedule]("payment schedule", s.PaymentSchedule),
		Payment:         NewCRUDHandler[models.Payment]("payment", s.Payment),
	}
}
