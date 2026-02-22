package handler

import (
	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/services"
)

type LeasingHandlers struct {
	LeasingProduct          ResourceHandler
	LeasingContract         ResourceHandler
	LeasingTask             ResourceHandler
	LeasingTaskAttribute    ResourceHandler
	LeasingContractDocument ResourceHandler
	Workflow                *LeasingWorkflowHandler
}

func NewLeasingHandlers(s services.LeasingServices) LeasingHandlers {
	return LeasingHandlers{
		LeasingProduct:          NewCRUDHandler[models.LeasingProduct]("leasing product", s.LeasingProduct),
		LeasingContract:         NewCRUDHandler[models.LeasingContract]("leasing contract", s.LeasingContract),
		LeasingTask:             NewCRUDHandler[models.LeasingTask]("leasing task", s.LeasingTask),
		LeasingTaskAttribute:    NewCRUDHandler[models.LeasingTaskAttribute]("leasing task attribute", s.LeasingTaskAttribute),
		LeasingContractDocument: NewCRUDHandler[models.LeasingContractDocument]("leasing contract document", s.LeasingContractDocument),
		Workflow:                NewLeasingWorkflowHandler(s.Workflow),
	}
}
