package models

import "time"

type MotorType struct {
	MotyID   int64   `gorm:"column:moty_id;primaryKey;autoIncrement"`
	MotyName string  `gorm:"column:moty_name;size:55;not null;uniqueIndex"`
	Motors   []Motor `gorm:"foreignKey:MotorMotyID;references:MotyID"`
}

func (MotorType) TableName() string { return "dealer.motor_types" }

type Motor struct {
	MotorID      int64        `gorm:"column:motor_id;primaryKey;autoIncrement"`
	Merk         string       `gorm:"column:merk;size:50;not null"`
	MotorType    string       `gorm:"column:motor_type;size:15;not null"`
	Tahun        int16        `gorm:"column:tahun;not null"`
	Warna        string       `gorm:"column:warna;size:30;not null"`
	NomorRangka  string       `gorm:"column:nomor_rangka;size:30;not null;uniqueIndex"`
	NomorMesin   string       `gorm:"column:nomor_mesin;size:30;not null;uniqueIndex"`
	CCMesin      string       `gorm:"column:cc_mesin;size:30;not null"`
	NomorPolisi  string       `gorm:"column:nomor_polisi;size:12;not null;uniqueIndex"`
	StatusUnit   string       `gorm:"column:status_unit;size:20;not null"`
	HargaOTR     float64      `gorm:"column:harga_otr;type:numeric(15,2);not null"`
	CreatedAt    time.Time    `gorm:"column:created_at;type:timestamptz;autoCreateTime"`
	MotorMotyID  int64        `gorm:"column:motor_moty_id;not null;index"`
	MotorTypeRef MotorType    `gorm:"foreignKey:MotorMotyID;references:MotyID"`
	MotorAssets  []MotorAsset `gorm:"foreignKey:MoasMotorID;references:MotorID"`
}

func (Motor) TableName() string { return "dealer.motors" }

type MotorAsset struct {
	MoasID      int64   `gorm:"column:moas_id;primaryKey;autoIncrement"`
	FileName    string  `gorm:"column:file_name;size:125;not null"`
	FileSize    float64 `gorm:"column:file_size;not null"`
	FileType    string  `gorm:"column:file_type;size:15;not null"`
	FileURL     string  `gorm:"column:file_url;size:125;not null"`
	MoasMotorID int64   `gorm:"column:moas_motor_id;not null;index"`
	Motor       Motor   `gorm:"foreignKey:MoasMotorID;references:MotorID"`
}

func (MotorAsset) TableName() string { return "dealer.motor_assets" }

type Customer struct {
	CustomerID       int64             `gorm:"column:customer_id;primaryKey;autoIncrement"`
	NIK              string            `gorm:"column:nik;size:16;not null;uniqueIndex"`
	NamaLengkap      string            `gorm:"column:nama_lengkap;size:100;not null"`
	TanggalLahir     time.Time         `gorm:"column:tanggal_lahir;type:date;not null"`
	NoHP             string            `gorm:"column:no_hp;size:15;not null;uniqueIndex"`
	Email            string            `gorm:"column:email;size:100;not null;uniqueIndex"`
	Pekerjaan        string            `gorm:"column:pekerjaan;size:80;not null"`
	Perusahaan       *string           `gorm:"column:perusahaan;size:120"`
	Salary           float64           `gorm:"column:salary;type:numeric(15,2);not null"`
	CreatedAt        time.Time         `gorm:"column:created_at;type:timestamptz;autoCreateTime"`
	UpdatedAt        time.Time         `gorm:"column:updated_at;type:timestamptz;autoUpdateTime"`
	LocationID       int64             `gorm:"column:location_id;not null;index"`
	Location         Location          `gorm:"foreignKey:LocationID;references:LocationID"`
	LeasingContracts []LeasingContract `gorm:"foreignKey:CustomerID;references:CustomerID"`
}

func (Customer) TableName() string { return "dealer.customer" }
