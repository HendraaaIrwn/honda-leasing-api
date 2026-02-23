package models

import "time"

type LeasingProduct struct {
	ProductID        int64             `gorm:"column:product_id;primaryKey;autoIncrement"`
	KodeProduk       string            `gorm:"column:kode_produk;size:20;not null;uniqueIndex"`
	NamaProduk       string            `gorm:"column:nama_produk;size:100;not null"`
	TenorBulan       int16             `gorm:"column:tenor_bulan;not null"`
	DPPersenMin      float64           `gorm:"column:dp_persen_min;type:numeric(5,2);not null"`
	DPPersenMax      float64           `gorm:"column:dp_persen_max;type:numeric(5,2);not null"`
	BungaFlat        float64           `gorm:"column:bunga_flat;type:numeric(5,2);not null"`
	AdminFee         float64           `gorm:"column:admin_fee;type:numeric(12,2);not null"`
	Asuransi         bool              `gorm:"column:asuransi;not null"`
	CreatedAt        time.Time         `gorm:"column:created_at;type:timestamptz;autoCreateTime"`
	LeasingContracts []LeasingContract `gorm:"foreignKey:ProductID;references:ProductID"`
}

func (LeasingProduct) TableName() string { return "leasing.leasing_product" }

type LeasingContract struct {
	ContractID        int64                     `gorm:"column:contract_id;primaryKey;autoIncrement"`
	ContractNumber    *string                   `gorm:"column:contract_number;size:30;uniqueIndex"`
	RequestDate       time.Time                 `gorm:"column:request_date;type:date;not null"`
	TanggalAkad       *time.Time                `gorm:"column:tanggal_akad;type:date"`
	TanggalMulaiCicil time.Time                 `gorm:"column:tanggal_mulai_cicil;type:date;not null"`
	TenorBulan        int16                     `gorm:"column:tenor_bulan;not null"`
	NilaiKendaraan    float64                   `gorm:"column:nilai_kendaraan;type:numeric(15,2);not null"`
	DPDibayar         float64                   `gorm:"column:dp_dibayar;type:numeric(15,2);not null"`
	PokokPinjaman     float64                   `gorm:"column:pokok_pinjaman;type:numeric(15,2);not null"`
	TotalPinjaman     float64                   `gorm:"column:total_pinjaman;type:numeric(15,2);not null"`
	CicilanPerBulan   float64                   `gorm:"column:cicilan_per_bulan;type:numeric(15,2);not null"`
	Status            string                    `gorm:"column:status;size:20;not null"`
	CreatedAt         time.Time                 `gorm:"column:created_at;type:timestamptz;autoCreateTime"`
	UpdatedAt         time.Time                 `gorm:"column:updated_at;type:timestamptz;autoUpdateTime"`
	CustomerID        int64                     `gorm:"column:customer_id;not null;index"`
	MotorID           int64                     `gorm:"column:motor_id;not null;index"`
	ProductID         int64                     `gorm:"column:product_id;not null;index"`
	Customer          Customer                  `gorm:"foreignKey:CustomerID;references:CustomerID"`
	Motor             Motor                     `gorm:"foreignKey:MotorID;references:MotorID"`
	Product           LeasingProduct            `gorm:"foreignKey:ProductID;references:ProductID"`
	LeasingTasks      []LeasingTask             `gorm:"foreignKey:ContractID;references:ContractID"`
	PaymentSchedules  []PaymentSchedule         `gorm:"foreignKey:ContractID;references:ContractID"`
	Payments          []Payment                 `gorm:"foreignKey:ContractID;references:ContractID"`
	ContractDocuments []LeasingContractDocument `gorm:"foreignKey:ContractID;references:ContractID"`
}

func (LeasingContract) TableName() string { return "leasing.leasing_contract" }

type LeasingTask struct {
	TaskID           int64                  `gorm:"column:task_id;primaryKey;autoIncrement"`
	TaskName         string                 `gorm:"column:task_name;size:85;not null"`
	StartDate        time.Time              `gorm:"column:startdate;type:date;not null"`
	EndDate          time.Time              `gorm:"column:enddate;type:date;not null"`
	ActualStartDate  *time.Time             `gorm:"column:actual_startdate;type:date"`
	ActualEndDate    *time.Time             `gorm:"column:actual_enddate;type:date"`
	SequenceNo       int                    `gorm:"column:sequence_no;not null"`
	Status           string                 `gorm:"column:status;size:15;not null"`
	ContractID       int64                  `gorm:"column:contract_id;not null;index"`
	RoleID           int64                  `gorm:"column:role_id;not null;index"`
	Contract         LeasingContract        `gorm:"foreignKey:ContractID;references:ContractID"`
	Role             Role                   `gorm:"foreignKey:RoleID;references:RoleID"`
	LeasingAttribute []LeasingTaskAttribute `gorm:"foreignKey:TasaLetaID;references:TaskID"`
}

func (LeasingTask) TableName() string { return "leasing.leasing_tasks" }

type LeasingTaskAttribute struct {
	TasaID     int64       `gorm:"column:tasa_id;primaryKey;autoIncrement"`
	TasaName   string      `gorm:"column:tasa_name;size:55;not null"`
	TasaValue  string      `gorm:"column:tasa_value;size:55;not null"`
	TasaStatus string      `gorm:"column:tasa_status;size:15;not null"`
	TasaLetaID int64       `gorm:"column:tasa_leta_id;not null;index"`
	Task       LeasingTask `gorm:"foreignKey:TasaLetaID;references:TaskID"`
}

func (LeasingTaskAttribute) TableName() string { return "leasing.leasing_tasks_attributes" }

type LeasingContractDocument struct {
	LocID      int64           `gorm:"column:loc_id;primaryKey;autoIncrement"`
	FileName   string          `gorm:"column:file_name;size:125;not null"`
	FileSize   float64         `gorm:"column:file_size;not null"`
	FileType   string          `gorm:"column:file_type;size:15;not null"`
	FileURL    string          `gorm:"column:file_url;size:125;not null"`
	ContractID int64           `gorm:"column:contract_id;not null;index"`
	Contract   LeasingContract `gorm:"foreignKey:ContractID;references:ContractID"`
}

func (LeasingContractDocument) TableName() string { return "leasing.leasing_contract_documents" }
