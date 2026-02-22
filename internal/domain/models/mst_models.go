package models

type Province struct {
	ProvID    int64       `gorm:"column:prov_id;primaryKey;autoIncrement"`
	ProvName  string      `gorm:"column:prov_name;size:85;not null;uniqueIndex"`
	Kabupaten []Kabupaten `gorm:"foreignKey:ProvID;references:ProvID"`
}

func (Province) TableName() string { return "mst.province" }

type Kabupaten struct {
	KabID     int64       `gorm:"column:kab_id;primaryKey;autoIncrement"`
	KabName   string      `gorm:"column:kab_name;size:85;not null"`
	ProvID    int64       `gorm:"column:prov_id;not null;index"`
	Province  Province    `gorm:"foreignKey:ProvID;references:ProvID"`
	Kecamatan []Kecamatan `gorm:"foreignKey:KabID;references:KabID"`
}

func (Kabupaten) TableName() string { return "mst.kabupaten" }

type Kecamatan struct {
	KecID     int64       `gorm:"column:kec_id;primaryKey;autoIncrement"`
	KecName   string      `gorm:"column:kec_name;size:85;not null"`
	KabID     int64       `gorm:"column:kab_id;not null;index"`
	Kabupaten Kabupaten   `gorm:"foreignKey:KabID;references:KabID"`
	Kelurahan []Kelurahan `gorm:"foreignKey:KecID;references:KecID"`
}

func (Kecamatan) TableName() string { return "mst.kecamatan" }

type Kelurahan struct {
	KelID     int64      `gorm:"column:kel_id;primaryKey;autoIncrement"`
	KelName   string     `gorm:"column:kel_name;size:85;not null"`
	KecID     int64      `gorm:"column:kec_id;not null;index"`
	Kecamatan Kecamatan  `gorm:"foreignKey:KecID;references:KecID"`
	Locations []Location `gorm:"foreignKey:KelID;references:KelID"`
}

func (Kelurahan) TableName() string { return "mst.kelurahan" }

type Location struct {
	LocationID    int64     `gorm:"column:location_id;primaryKey;autoIncrement"`
	StreetAddress string    `gorm:"column:street_address;type:text"`
	PostalCode    string    `gorm:"column:postal_code;size:10"`
	Longitude     string    `gorm:"column:longitude;type:text"`
	Latitude      string    `gorm:"column:latitude;type:text"`
	KelID         int64     `gorm:"column:kel_id;not null;index"`
	Kelurahan     Kelurahan `gorm:"foreignKey:KelID;references:KelID"`
}

func (Location) TableName() string { return "mst.locations" }

type TemplateTask struct {
	TetaID     int64                   `gorm:"column:teta_id;primaryKey;autoIncrement"`
	TetaName   string                  `gorm:"column:teta_name;size:85;not null"`
	TetaRoleID int64                   `gorm:"column:teta_role_id;not null;index"`
	Role       Role                    `gorm:"foreignKey:TetaRoleID;references:RoleID"`
	Attributes []TemplateTaskAttribute `gorm:"foreignKey:TetatTetaID;references:TetaID"`
}

func (TemplateTask) TableName() string { return "mst.template_tasks" }

type TemplateTaskAttribute struct {
	TetatID      int64        `gorm:"column:tetat_id;primaryKey;autoIncrement"`
	TetatName    string       `gorm:"column:tetat_name;size:85;not null"`
	TetatTetaID  int64        `gorm:"column:tetat_teta_id;not null;index"`
	TemplateTask TemplateTask `gorm:"foreignKey:TetatTetaID;references:TetaID"`
}

func (TemplateTaskAttribute) TableName() string { return "mst.template_task_attributes" }
