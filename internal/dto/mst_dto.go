package dto

type ProvinceDTO struct {
	ProvID   int64  `json:"prov_id"`
	ProvName string `json:"prov_name"`
}

type KabupatenDTO struct {
	KabID   int64  `json:"kab_id"`
	KabName string `json:"kab_name"`
	ProvID  int64  `json:"prov_id"`
}

type KecamatanDTO struct {
	KecID   int64  `json:"kec_id"`
	KecName string `json:"kec_name"`
	KabID   int64  `json:"kab_id"`
}

type KelurahanDTO struct {
	KelID   int64  `json:"kel_id"`
	KelName string `json:"kel_name"`
	KecID   int64  `json:"kec_id"`
}

type LocationDTO struct {
	LocationID    int64  `json:"location_id"`
	StreetAddress string `json:"street_address"`
	PostalCode    string `json:"postal_code"`
	Longitude     string `json:"longitude"`
	Latitude      string `json:"latitude"`
	KelID         int64  `json:"kel_id"`
}

type TemplateTaskDTO struct {
	TetaID     int64  `json:"teta_id"`
	TetaName   string `json:"teta_name"`
	TetaRoleID int64  `json:"teta_role_id"`
}

type TemplateTaskAttributeDTO struct {
	TetatID     int64  `json:"tetat_id"`
	TetatName   string `json:"tetat_name"`
	TetatTetaID int64  `json:"tetat_teta_id"`
}
