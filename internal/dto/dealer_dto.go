package dto

import "time"

type MotorTypeDTO struct {
	MotyID   int64  `json:"moty_id"`
	MotyName string `json:"moty_name"`
}

type MotorDTO struct {
	MotorID     int64     `json:"motor_id"`
	Merk        string    `json:"merk"`
	MotorType   string    `json:"motor_type"`
	Tahun       int16     `json:"tahun"`
	Warna       string    `json:"warna"`
	NomorRangka string    `json:"nomor_rangka"`
	NomorMesin  string    `json:"nomor_mesin"`
	CCMesin     string    `json:"cc_mesin"`
	NomorPolisi string    `json:"nomor_polisi"`
	StatusUnit  string    `json:"status_unit"`
	HargaOTR    float64   `json:"harga_otr"`
	CreatedAt   time.Time `json:"created_at"`
	MotorMotyID int64     `json:"motor_moty_id"`
}

type MotorAssetDTO struct {
	MoasID      int64   `json:"moas_id"`
	FileName    string  `json:"file_name"`
	FileSize    float64 `json:"file_size"`
	FileType    string  `json:"file_type"`
	FileURL     string  `json:"file_url"`
	MoasMotorID int64   `json:"moas_motor_id"`
}

type CustomerDTO struct {
	CustomerID   int64     `json:"customer_id"`
	NIK          string    `json:"nik"`
	NamaLengkap  string    `json:"nama_lengkap"`
	TanggalLahir time.Time `json:"tanggal_lahir"`
	NoHP         string    `json:"no_hp"`
	Email        string    `json:"email"`
	Pekerjaan    string    `json:"pekerjaan"`
	Perusahaan   *string   `json:"perusahaan"`
	Salary       float64   `json:"salary"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LocationID   int64     `json:"location_id"`
}
