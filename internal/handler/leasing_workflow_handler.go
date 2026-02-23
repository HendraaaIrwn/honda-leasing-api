package handler

import (
	"time"

	errs "github.com/HendraaaIrwn/honda-leasing-api/internal/errors"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/response"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/services"
	"github.com/gin-gonic/gin"
)

type LeasingWorkflowHandler struct {
	service services.LeasingWorkflowService
}

func NewLeasingWorkflowHandler(service services.LeasingWorkflowService) *LeasingWorkflowHandler {
	if service == nil {
		return nil
	}

	return &LeasingWorkflowHandler{service: service}
}

func (h *LeasingWorkflowHandler) RegisterRoutes(group *gin.RouterGroup) {
	workflow := group.Group("/workflow")
	workflow.POST("/submit-application", h.SubmitApplication)
	workflow.POST("/auto-scoring", h.ProcessAutoScoring)
	workflow.POST("/survey", h.ProcessSurveyResult)
	workflow.POST("/final-approval", h.ProcessFinalApproval)
	workflow.POST("/akad", h.ExecuteAkad)
	workflow.POST("/initial-payment", h.RecordInitialPayment)
	workflow.POST("/dealer-fulfillment", h.ProcessDealerFulfillment)
	workflow.POST("/delivery", h.CompleteDelivery)
}

type contractDocumentRequest struct {
	FileName string  `json:"file_name"`
	FileSize float64 `json:"file_size"`
	FileType string  `json:"file_type"`
	FileURL  string  `json:"file_url"`
}

type submitApplicationRequest struct {
	CustomerID  int64                     `json:"customer_id"`
	MotorID     int64                     `json:"motor_id"`
	ProductID   int64                     `json:"product_id"`
	DPDibayar   float64                   `json:"dp_dibayar"`
	TenorBulan  int16                     `json:"tenor_bulan"`
	RequestDate *time.Time                `json:"request_date"`
	Documents   []contractDocumentRequest `json:"documents"`
}

type autoScoringDecisionRequest struct {
	ContractID        int64  `json:"contract_id"`
	AutoApproved      bool   `json:"auto_approved"`
	ManualReviewReady bool   `json:"manual_review_ready"`
	ManualApproved    bool   `json:"manual_approved"`
	Note              string `json:"note"`
}

type surveyDecisionRequest struct {
	ContractID   int64   `json:"contract_id"`
	Decision     string  `json:"decision"`
	AdditionalDP float64 `json:"additional_dp"`
	Note         string  `json:"note"`
}

type finalApprovalRequest struct {
	ContractID int64  `json:"contract_id"`
	Approved   bool   `json:"approved"`
	Note       string `json:"note"`
}

type akadRequest struct {
	ContractID           int64      `json:"contract_id"`
	ContractNumber       string     `json:"contract_number"`
	AkadDate             *time.Time `json:"akad_date"`
	TanggalMulaiCicil    *time.Time `json:"tanggal_mulai_cicil"`
	GenerateContractCode bool       `json:"generate_contract_code"`
}

type initialPaymentRequest struct {
	ContractID       int64      `json:"contract_id"`
	NomorBukti       string     `json:"nomor_bukti"`
	JumlahBayar      float64    `json:"jumlah_bayar"`
	TanggalBayar     *time.Time `json:"tanggal_bayar"`
	MetodePembayaran string     `json:"metode_pembayaran"`
	Provider         string     `json:"provider"`
}

type dealerFulfillmentRequest struct {
	ContractID          int64  `json:"contract_id"`
	UnitReadyStock      bool   `json:"unit_ready_stock"`
	EstimatedIndentWeek int    `json:"estimated_indent_week"`
	Note                string `json:"note"`
}

type deliveryCompletionRequest struct {
	ContractID         int64                     `json:"contract_id"`
	DeliveryDate       *time.Time                `json:"delivery_date"`
	CustomerReceived   bool                      `json:"customer_received"`
	DocumentHandover   bool                      `json:"document_handover"`
	HandoverNote       string                    `json:"handover_note"`
	TanggalMulaiCicil  *time.Time                `json:"tanggal_mulai_cicil"`
	ContractDocUploads []contractDocumentRequest `json:"contract_doc_uploads"`
}

func (h *LeasingWorkflowHandler) SubmitApplication(c *gin.Context) {
	var req submitApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.ErrInvalidInput)
		return
	}

	result, err := h.service.SubmitApplication(c.Request.Context(), services.SubmitApplicationInput{
		CustomerID:  req.CustomerID,
		MotorID:     req.MotorID,
		ProductID:   req.ProductID,
		DPDibayar:   req.DPDibayar,
		TenorBulan:  req.TenorBulan,
		RequestDate: req.RequestDate,
		Documents:   mapContractDocs(req.Documents),
	})
	if err != nil {
		respondError(c, err)
		return
	}

	response.Created(c, "application submitted", result)
}

func (h *LeasingWorkflowHandler) ProcessAutoScoring(c *gin.Context) {
	var req autoScoringDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.ErrInvalidInput)
		return
	}

	err := h.service.ProcessAutoScoring(c.Request.Context(), services.AutoScoringDecisionInput{
		ContractID:        req.ContractID,
		AutoApproved:      req.AutoApproved,
		ManualReviewReady: req.ManualReviewReady,
		ManualApproved:    req.ManualApproved,
		Note:              req.Note,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, "auto scoring processed", gin.H{"contract_id": req.ContractID})
}

func (h *LeasingWorkflowHandler) ProcessSurveyResult(c *gin.Context) {
	var req surveyDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.ErrInvalidInput)
		return
	}

	err := h.service.ProcessSurveyResult(c.Request.Context(), services.SurveyDecisionInput{
		ContractID:   req.ContractID,
		Decision:     services.SurveyDecision(req.Decision),
		AdditionalDP: req.AdditionalDP,
		Note:         req.Note,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, "survey result processed", gin.H{"contract_id": req.ContractID})
}

func (h *LeasingWorkflowHandler) ProcessFinalApproval(c *gin.Context) {
	var req finalApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.ErrInvalidInput)
		return
	}

	err := h.service.ProcessFinalApproval(c.Request.Context(), services.FinalApprovalInput{
		ContractID: req.ContractID,
		Approved:   req.Approved,
		Note:       req.Note,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, "final approval processed", gin.H{"contract_id": req.ContractID})
}

func (h *LeasingWorkflowHandler) ExecuteAkad(c *gin.Context) {
	var req akadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.ErrInvalidInput)
		return
	}

	akadDate := time.Time{}
	if req.AkadDate != nil {
		akadDate = *req.AkadDate
	}

	err := h.service.ExecuteAkad(c.Request.Context(), services.AkadInput{
		ContractID:           req.ContractID,
		ContractNumber:       req.ContractNumber,
		AkadDate:             akadDate,
		TanggalMulaiCicil:    req.TanggalMulaiCicil,
		GenerateContractCode: req.GenerateContractCode,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, "akad processed", gin.H{"contract_id": req.ContractID})
}

func (h *LeasingWorkflowHandler) RecordInitialPayment(c *gin.Context) {
	var req initialPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.ErrInvalidInput)
		return
	}

	tanggalBayar := time.Time{}
	if req.TanggalBayar != nil {
		tanggalBayar = *req.TanggalBayar
	}

	err := h.service.RecordInitialPayment(c.Request.Context(), services.InitialPaymentInput{
		ContractID:       req.ContractID,
		NomorBukti:       req.NomorBukti,
		JumlahBayar:      req.JumlahBayar,
		TanggalBayar:     tanggalBayar,
		MetodePembayaran: req.MetodePembayaran,
		Provider:         req.Provider,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, "initial payment recorded", gin.H{"contract_id": req.ContractID})
}

func (h *LeasingWorkflowHandler) ProcessDealerFulfillment(c *gin.Context) {
	var req dealerFulfillmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.ErrInvalidInput)
		return
	}

	err := h.service.ProcessDealerFulfillment(c.Request.Context(), services.DealerFulfillmentInput{
		ContractID:          req.ContractID,
		UnitReadyStock:      req.UnitReadyStock,
		EstimatedIndentWeek: req.EstimatedIndentWeek,
		Note:                req.Note,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, "dealer fulfillment processed", gin.H{"contract_id": req.ContractID})
}

func (h *LeasingWorkflowHandler) CompleteDelivery(c *gin.Context) {
	var req deliveryCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errs.ErrInvalidInput)
		return
	}

	deliveryDate := time.Time{}
	if req.DeliveryDate != nil {
		deliveryDate = *req.DeliveryDate
	}

	err := h.service.CompleteDelivery(c.Request.Context(), services.DeliveryCompletionInput{
		ContractID:         req.ContractID,
		DeliveryDate:       deliveryDate,
		CustomerReceived:   req.CustomerReceived,
		DocumentHandover:   req.DocumentHandover,
		HandoverNote:       req.HandoverNote,
		TanggalMulaiCicil:  req.TanggalMulaiCicil,
		ContractDocUploads: mapContractDocs(req.ContractDocUploads),
	})
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, "delivery completed", gin.H{"contract_id": req.ContractID})
}

func mapContractDocs(input []contractDocumentRequest) []services.ContractDocumentInput {
	if len(input) == 0 {
		return nil
	}

	result := make([]services.ContractDocumentInput, 0, len(input))
	for _, doc := range input {
		result = append(result, services.ContractDocumentInput{
			FileName: doc.FileName,
			FileSize: doc.FileSize,
			FileType: doc.FileType,
			FileURL:  doc.FileURL,
		})
	}

	return result
}
