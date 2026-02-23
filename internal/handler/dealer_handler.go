package handler

import (
	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/services"
)

type DealerHandlers struct {
	MotorType  ResourceHandler
	Motor      ResourceHandler
	MotorAsset ResourceHandler
	Customer   ResourceHandler
}

func NewDealerHandlers(s services.DealerServices) DealerHandlers {
	return DealerHandlers{
		MotorType:  NewCRUDHandler[models.MotorType]("motor type", s.MotorType),
		Motor:      NewCRUDHandler[models.Motor]("motor", s.Motor),
		MotorAsset: NewCRUDHandler[models.MotorAsset]("motor asset", s.MotorAsset),
		Customer:   NewCRUDHandler[models.Customer]("customer", s.Customer),
	}
}
