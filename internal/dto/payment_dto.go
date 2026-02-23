package dto

import "time"

type PaymentScheduleDTO struct {
	ScheduleID       int64      `json:"schedule_id"`
	AngsuranKe       int16      `json:"angsuran_ke"`
	JatuhTempo       time.Time  `json:"jatuh_tempo"`
	Pokok            float64    `json:"pokok"`
	Margin           float64    `json:"margin"`
	TotalTagihan     float64    `json:"total_tagihan"`
	StatusPembayaran string     `json:"status_pembayaran"`
	TanggalBayar     *time.Time `json:"tanggal_bayar"`
	CreatedAt        time.Time  `json:"created_at"`
	ContractID       int64      `json:"contract_id"`
}

type PaymentDTO struct {
	PaymentID        int64     `json:"payment_id"`
	NomorBukti       string    `json:"nomor_bukti"`
	JumlahBayar      float64   `json:"jumlah_bayar"`
	TanggalBayar     time.Time `json:"tanggal_bayar"`
	MetodePembayaran string    `json:"metode_pembayaran"`
	Provider         string    `json:"provider"`
	CreatedAt        time.Time `json:"created_at"`
	ContractID       int64     `json:"contract_id"`
	ScheduleID       *int64    `json:"schedule_id"`
}
