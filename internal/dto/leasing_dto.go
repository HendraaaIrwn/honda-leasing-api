package dto

import "time"

type LeasingProductDTO struct {
	ProductID   int64     `json:"product_id"`
	KodeProduk  string    `json:"kode_produk"`
	NamaProduk  string    `json:"nama_produk"`
	TenorBulan  int16     `json:"tenor_bulan"`
	DPPersenMin float64   `json:"dp_persen_min"`
	DPPersenMax float64   `json:"dp_persen_max"`
	BungaFlat   float64   `json:"bunga_flat"`
	AdminFee    float64   `json:"admin_fee"`
	Asuransi    bool      `json:"asuransi"`
	CreatedAt   time.Time `json:"created_at"`
}

type LeasingContractDTO struct {
	ContractID        int64      `json:"contract_id"`
	ContractNumber    *string    `json:"contract_number"`
	RequestDate       time.Time  `json:"request_date"`
	TanggalAkad       *time.Time `json:"tanggal_akad"`
	TanggalMulaiCicil time.Time  `json:"tanggal_mulai_cicil"`
	TenorBulan        int16      `json:"tenor_bulan"`
	NilaiKendaraan    float64    `json:"nilai_kendaraan"`
	DPDibayar         float64    `json:"dp_dibayar"`
	PokokPinjaman     float64    `json:"pokok_pinjaman"`
	TotalPinjaman     float64    `json:"total_pinjaman"`
	CicilanPerBulan   float64    `json:"cicilan_per_bulan"`
	Status            string     `json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CustomerID        int64      `json:"customer_id"`
	MotorID           int64      `json:"motor_id"`
	ProductID         int64      `json:"product_id"`
}

type LeasingTaskDTO struct {
	TaskID          int64      `json:"task_id"`
	TaskName        string     `json:"task_name"`
	StartDate       time.Time  `json:"startdate"`
	EndDate         time.Time  `json:"enddate"`
	ActualStartDate *time.Time `json:"actual_startdate"`
	ActualEndDate   *time.Time `json:"actual_enddate"`
	SequenceNo      int        `json:"sequence_no"`
	Status          string     `json:"status"`
	ContractID      int64      `json:"contract_id"`
	RoleID          int64      `json:"role_id"`
}

type LeasingTaskAttributeDTO struct {
	TasaID     int64  `json:"tasa_id"`
	TasaName   string `json:"tasa_name"`
	TasaValue  string `json:"tasa_value"`
	TasaStatus string `json:"tasa_status"`
	TasaLetaID int64  `json:"tasa_leta_id"`
}

type LeasingContractDocumentDTO struct {
	LocID      int64   `json:"loc_id"`
	FileName   string  `json:"file_name"`
	FileSize   float64 `json:"file_size"`
	FileType   string  `json:"file_type"`
	FileURL    string  `json:"file_url"`
	ContractID int64   `json:"contract_id"`
}
