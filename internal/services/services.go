package services

import (
	"github.com/HendraaaIrwn/honda-leasing-api/internal/repository"
	"github.com/HendraaaIrwn/honda-leasing-api/pkg/database"
)

type AccountServices struct {
	OAuthProvider     OAuthProviderService
	User              UserService
	UserOAuthProvider UserOAuthProviderService
	Role              RoleService
	UserRole          UserRoleService
	Permission        PermissionService
	RolePermission    RolePermissionService
}

type MSTServices struct {
	Province              ProvinceService
	Kabupaten             KabupatenService
	Kecamatan             KecamatanService
	Kelurahan             KelurahanService
	Location              LocationService
	TemplateTask          TemplateTaskService
	TemplateTaskAttribute TemplateTaskAttributeService
}

type DealerServices struct {
	MotorType  MotorTypeService
	Motor      MotorService
	MotorAsset MotorAssetService
	Customer   CustomerService
}

type LeasingServices struct {
	LeasingProduct          LeasingProductService
	LeasingContract         LeasingContractService
	LeasingTask             LeasingTaskService
	LeasingTaskAttribute    LeasingTaskAttributeService
	LeasingContractDocument LeasingContractDocumentService
	Workflow                LeasingWorkflowService
}

type PaymentServices struct {
	PaymentSchedule PaymentScheduleService
	Payment         PaymentService
}

// Services is the domain service registry.
type Services struct {
	Account AccountServices
	MST     MSTServices
	Dealer  DealerServices
	Leasing LeasingServices
	Payment PaymentServices
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Account: AccountServices{
			OAuthProvider:     NewOAuthProviderService(repos.Account.OAuthProvider),
			User:              NewUserService(repos.Account.User),
			UserOAuthProvider: NewUserOAuthProviderService(repos.Account.UserOAuthProvider),
			Role:              NewRoleService(repos.Account.Role),
			UserRole:          NewUserRoleService(repos.Account.UserRole),
			Permission:        NewPermissionService(repos.Account.Permission),
			RolePermission:    NewRolePermissionService(repos.Account.RolePermission),
		},
		MST: MSTServices{
			Province:              NewProvinceService(repos.MST.Province),
			Kabupaten:             NewKabupatenService(repos.MST.Kabupaten),
			Kecamatan:             NewKecamatanService(repos.MST.Kecamatan),
			Kelurahan:             NewKelurahanService(repos.MST.Kelurahan),
			Location:              NewLocationService(repos.MST.Location),
			TemplateTask:          NewTemplateTaskService(repos.MST.TemplateTask),
			TemplateTaskAttribute: NewTemplateTaskAttributeService(repos.MST.TemplateTaskAttribute),
		},
		Dealer: DealerServices{
			MotorType:  NewMotorTypeService(repos.Dealer.MotorType),
			Motor:      NewMotorService(repos.Dealer.Motor),
			MotorAsset: NewMotorAssetService(repos.Dealer.MotorAsset),
			Customer:   NewCustomerService(repos.Dealer.Customer),
		},
		Leasing: LeasingServices{
			LeasingProduct:          NewLeasingProductService(repos.Leasing.LeasingProduct),
			LeasingContract:         NewLeasingContractService(repos.Leasing.LeasingContract),
			LeasingTask:             NewLeasingTaskService(repos.Leasing.LeasingTask),
			LeasingTaskAttribute:    NewLeasingTaskAttributeService(repos.Leasing.LeasingTaskAttribute),
			LeasingContractDocument: NewLeasingContractDocumentService(repos.Leasing.LeasingContractDocument),
			Workflow:                NewLeasingWorkflowService(repos.DB()),
		},
		Payment: PaymentServices{
			PaymentSchedule: NewPaymentScheduleService(repos.Payment.PaymentSchedule),
			Payment:         NewPaymentService(repos.Payment.Payment),
		},
	}
}

func NewServicesFromDatabase(db *database.Database) *Services {
	repos := repository.NewRepositoriesFromDatabase(db)
	return NewServices(repos)
}
