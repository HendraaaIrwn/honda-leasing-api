package handler

import "github.com/HendraaaIrwn/honda-leasing-api/internal/services"

// Handlers is a registry for all domain handlers.
type Handlers struct {
	Account AccountHandlers
	MST     MSTHandlers
	Dealer  DealerHandlers
	Leasing LeasingHandlers
	Payment PaymentHandlers
}

func NewHandlers(s *services.Services) *Handlers {
	return &Handlers{
		Account: NewAccountHandlers(s.Account),
		MST:     NewMSTHandlers(s.MST),
		Dealer:  NewDealerHandlers(s.Dealer),
		Leasing: NewLeasingHandlers(s.Leasing),
		Payment: NewPaymentHandlers(s.Payment),
	}
}
