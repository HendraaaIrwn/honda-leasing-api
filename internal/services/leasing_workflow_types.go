package services

import "time"

const (
	ContractStatusDraft    = "draft"
	ContractStatusApproved = "approved"
	ContractStatusActive   = "active"
	ContractStatusCanceled = "canceled"

	MotorStatusReady  = "ready"
	MotorStatusBooked = "booked"
	MotorStatusLeased = "leased"

	TaskStatusInProgress = "inprogress"
	TaskStatusCompleted  = "completed"
	TaskStatusCancelled  = "cancelled"

	TaskAttrStatusPending   = "pending"
	TaskAttrStatusCompleted = "completed"
	TaskAttrStatusCancelled = "cancelled"
)

type SurveyDecision string

const (
	SurveyDecisionApprove             SurveyDecision = "approve"
	SurveyDecisionReject              SurveyDecision = "reject"
	SurveyDecisionRequestAdditionalDP SurveyDecision = "request_additional_dp"
)

type ContractDocumentInput struct {
	FileName string
	FileSize float64
	FileType string
	FileURL  string
}

type SubmitApplicationInput struct {
	CustomerID int64
	MotorID    int64
	ProductID  int64
	DPDibayar  float64
	TenorBulan int16

	RequestDate *time.Time
	Documents   []ContractDocumentInput
}

type AutoScoringDecisionInput struct {
	ContractID        int64
	AutoApproved      bool
	ManualReviewReady bool
	ManualApproved    bool
	Note              string
}

type SurveyDecisionInput struct {
	ContractID   int64
	Decision     SurveyDecision
	AdditionalDP float64
	Note         string
}

type FinalApprovalInput struct {
	ContractID int64
	Approved   bool
	Note       string
}

type AkadInput struct {
	ContractID           int64
	ContractNumber       string
	AkadDate             time.Time
	TanggalMulaiCicil    *time.Time
	GenerateContractCode bool
}

type InitialPaymentInput struct {
	ContractID       int64
	NomorBukti       string
	JumlahBayar      float64
	TanggalBayar     time.Time
	MetodePembayaran string
	Provider         string
}

type DealerFulfillmentInput struct {
	ContractID          int64
	UnitReadyStock      bool
	EstimatedIndentWeek int
	Note                string
}

type DeliveryCompletionInput struct {
	ContractID         int64
	DeliveryDate       time.Time
	CustomerReceived   bool
	DocumentHandover   bool
	HandoverNote       string
	TanggalMulaiCicil  *time.Time
	ContractDocUploads []ContractDocumentInput
}
