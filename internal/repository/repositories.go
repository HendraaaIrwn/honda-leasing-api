package repository

import (
	"github.com/HendraaaIrwn/honda-leasing-api/pkg/database"
	"gorm.io/gorm"
)

type AccountRepositories struct {
	OAuthProvider     OAuthProviderRepository
	User              UserRepository
	UserOAuthProvider UserOAuthProviderRepository
	Role              RoleRepository
	UserRole          UserRoleRepository
	Permission        PermissionRepository
	RolePermission    RolePermissionRepository
}

type MSTRepositories struct {
	Province              ProvinceRepository
	Kabupaten             KabupatenRepository
	Kecamatan             KecamatanRepository
	Kelurahan             KelurahanRepository
	Location              LocationRepository
	TemplateTask          TemplateTaskRepository
	TemplateTaskAttribute TemplateTaskAttributeRepository
}

type DealerRepositories struct {
	MotorType  MotorTypeRepository
	Motor      MotorRepository
	MotorAsset MotorAssetRepository
	Customer   CustomerRepository
}

type LeasingRepositories struct {
	LeasingProduct          LeasingProductRepository
	LeasingContract         LeasingContractRepository
	LeasingTask             LeasingTaskRepository
	LeasingTaskAttribute    LeasingTaskAttributeRepository
	LeasingContractDocument LeasingContractDocumentRepository
}

type PaymentRepositories struct {
	PaymentSchedule PaymentScheduleRepository
	Payment         PaymentRepository
}

// Repositories is the domain repository registry.
type Repositories struct {
	db      *gorm.DB
	Account AccountRepositories
	MST     MSTRepositories
	Dealer  DealerRepositories
	Leasing LeasingRepositories
	Payment PaymentRepositories
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		db: db,
		Account: AccountRepositories{
			OAuthProvider:     NewOAuthProviderRepository(db),
			User:              NewUserRepository(db),
			UserOAuthProvider: NewUserOAuthProviderRepository(db),
			Role:              NewRoleRepository(db),
			UserRole:          NewUserRoleRepository(db),
			Permission:        NewPermissionRepository(db),
			RolePermission:    NewRolePermissionRepository(db),
		},
		MST: MSTRepositories{
			Province:              NewProvinceRepository(db),
			Kabupaten:             NewKabupatenRepository(db),
			Kecamatan:             NewKecamatanRepository(db),
			Kelurahan:             NewKelurahanRepository(db),
			Location:              NewLocationRepository(db),
			TemplateTask:          NewTemplateTaskRepository(db),
			TemplateTaskAttribute: NewTemplateTaskAttributeRepository(db),
		},
		Dealer: DealerRepositories{
			MotorType:  NewMotorTypeRepository(db),
			Motor:      NewMotorRepository(db),
			MotorAsset: NewMotorAssetRepository(db),
			Customer:   NewCustomerRepository(db),
		},
		Leasing: LeasingRepositories{
			LeasingProduct:          NewLeasingProductRepository(db),
			LeasingContract:         NewLeasingContractRepository(db),
			LeasingTask:             NewLeasingTaskRepository(db),
			LeasingTaskAttribute:    NewLeasingTaskAttributeRepository(db),
			LeasingContractDocument: NewLeasingContractDocumentRepository(db),
		},
		Payment: PaymentRepositories{
			PaymentSchedule: NewPaymentScheduleRepository(db),
			Payment:         NewPaymentRepository(db),
		},
	}
}

func (r *Repositories) DB() *gorm.DB {
	return r.db
}

func NewRepositoriesFromDatabase(db *database.Database) *Repositories {
	return NewRepositories(database.GetDB(db))
}
