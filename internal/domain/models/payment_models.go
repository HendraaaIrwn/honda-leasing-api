package models

import "time"

type PaymentSchedule struct {
	ScheduleID       int64           `gorm:"column:schedule_id;primaryKey;autoIncrement"`
	AngsuranKe       int16           `gorm:"column:angsuran_ke;not null"`
	JatuhTempo       time.Time       `gorm:"column:jatuh_tempo;type:date;not null"`
	Pokok            float64         `gorm:"column:pokok;type:numeric(15,2);not null"`
	Margin           float64         `gorm:"column:margin;type:numeric(15,2);not null"`
	TotalTagihan     float64         `gorm:"column:total_tagihan;type:numeric(15,2);not null"`
	StatusPembayaran string          `gorm:"column:status_pembayaran;size:20;not null"`
	TanggalBayar     *time.Time      `gorm:"column:tanggal_bayar;type:date"`
	CreatedAt        time.Time       `gorm:"column:created_at;type:timestamptz;autoCreateTime"`
	ContractID       int64           `gorm:"column:contract_id;not null;index"`
	Contract         LeasingContract `gorm:"foreignKey:ContractID;references:ContractID"`
	Payments         []Payment       `gorm:"foreignKey:ScheduleID;references:ScheduleID"`
}

func (PaymentSchedule) TableName() string { return "payment.payment_schedule" }

type Payment struct {
	PaymentID        int64            `gorm:"column:payment_id;primaryKey;autoIncrement"`
	NomorBukti       string           `gorm:"column:nomor_bukti;size:40;not null;uniqueIndex"`
	JumlahBayar      float64          `gorm:"column:jumlah_bayar;type:numeric(15,2);not null"`
	TanggalBayar     time.Time        `gorm:"column:tanggal_bayar;type:date;not null"`
	MetodePembayaran string           `gorm:"column:metode_pembayaran;size:30;not null"`
	Provider         string           `gorm:"column:provider;size:50;not null"`
	CreatedAt        time.Time        `gorm:"column:created_at;type:timestamptz;autoCreateTime"`
	ContractID       int64            `gorm:"column:contract_id;not null;index"`
	ScheduleID       *int64           `gorm:"column:schedule_id;index"`
	Contract         LeasingContract  `gorm:"foreignKey:ContractID;references:ContractID"`
	Schedule         *PaymentSchedule `gorm:"foreignKey:ScheduleID;references:ScheduleID"`
}

func (Payment) TableName() string { return "payment.payments" }
