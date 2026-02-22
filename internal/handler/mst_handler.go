package handler

import (
	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/services"
)

type MSTHandlers struct {
	Province              ResourceHandler
	Kabupaten             ResourceHandler
	Kecamatan             ResourceHandler
	Kelurahan             ResourceHandler
	Location              ResourceHandler
	TemplateTask          ResourceHandler
	TemplateTaskAttribute ResourceHandler
}

func NewMSTHandlers(s services.MSTServices) MSTHandlers {
	return MSTHandlers{
		Province:              NewCRUDHandler[models.Province]("province", s.Province),
		Kabupaten:             NewCRUDHandler[models.Kabupaten]("kabupaten", s.Kabupaten),
		Kecamatan:             NewCRUDHandler[models.Kecamatan]("kecamatan", s.Kecamatan),
		Kelurahan:             NewCRUDHandler[models.Kelurahan]("kelurahan", s.Kelurahan),
		Location:              NewCRUDHandler[models.Location]("location", s.Location),
		TemplateTask:          NewCRUDHandler[models.TemplateTask]("template task", s.TemplateTask),
		TemplateTaskAttribute: NewCRUDHandler[models.TemplateTaskAttribute]("template task attribute", s.TemplateTaskAttribute),
	}
}
